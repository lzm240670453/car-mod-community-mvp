<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">帖子后审</h1>
        <div class="page-desc">先发后审与高风险关键词待审内容统一处理。</div>
      </div>
      <el-switch v-model="pendingOnly" active-text="只看待审核" @change="load" />
    </div>

    <div class="toolbar">
      <el-select v-model="filters.type" clearable placeholder="帖子类型" style="width: 160px">
        <el-option v-for="item in POST_TYPES" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-input-number v-model="filters.userId" :min="1" controls-position="right" placeholder="用户 ID" style="width: 160px" />
      <el-button type="primary" :icon="Search" @click="load">查询</el-button>
    </div>

    <div class="table-shell">
      <el-table v-loading="loading" :data="items" row-key="id">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="内容" min-width="340">
          <template #default="{ row }">
            <div class="text-strong">{{ row.title }}</div>
            <div class="content-cell text-muted">{{ row.content }}</div>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="120">
          <template #default="{ row }">{{ labelOf(POST_TYPES, row.type) }}</template>
        </el-table-column>
        <el-table-column prop="userId" label="用户" width="90" />
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ labelOf(CONTENT_STATUSES, row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="互动" width="160">
          <template #default="{ row }">{{ row.likeCount }} 赞 / {{ row.commentCount }} 评 / {{ row.favoriteCount }} 藏</template>
        </el-table-column>
        <el-table-column label="发布时间" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="success" plain @click="audit(row.id, 'approve')">通过</el-button>
            <el-button size="small" type="danger" plain @click="audit(row.id, 'hide')">隐藏</el-button>
            <el-button size="small" plain @click="audit(row.id, 'restore')">恢复</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        layout="total, sizes, prev, pager, next"
        :total="total"
        :page-sizes="[20, 50, 100]"
        style="margin-top: 14px; justify-content: flex-end"
        @change="load"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessageBox } from "element-plus";
import { Search } from "@element-plus/icons-vue";
import { approvePost, hidePost, listPendingReviewPosts, listPostsForReview, restorePost } from "../api/admin";
import type { Post } from "../types";
import { CONTENT_STATUSES, formatDate, labelOf, POST_TYPES, statusTagType } from "../utils";

const loading = ref(false);
const pendingOnly = ref(false);
const items = ref<Post[]>([]);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);
const filters = reactive<{ type?: number; userId?: number }>({});

async function load() {
  loading.value = true;
  try {
    const params = { page: page.value, pageSize: pageSize.value, type: filters.type, userId: filters.userId };
    const result = pendingOnly.value ? await listPendingReviewPosts(params) : await listPostsForReview(params);
    items.value = result.items;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

async function audit(id: number, action: "approve" | "hide" | "restore") {
  const title = action === "approve" ? "审核通过" : action === "hide" ? "隐藏帖子" : "恢复帖子";
  const { value } = await ElMessageBox.prompt("请输入审核备注", title, {
    inputPlaceholder: "备注",
    confirmButtonText: "确认",
    cancelButtonText: "取消"
  });
  if (action === "approve") {
    await approvePost(id, value || "");
  } else if (action === "hide") {
    await hidePost(id, value || "");
  } else {
    await restorePost(id, value || "");
  }
  await load();
}

onMounted(load);
</script>
