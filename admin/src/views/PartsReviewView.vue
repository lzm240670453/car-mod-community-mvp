<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">二手件后审</h1>
        <div class="page-desc">处理出售、求购信息中的待审和已展示内容。</div>
      </div>
      <el-switch v-model="pendingOnly" active-text="只看待审核" @change="load" />
    </div>

    <div class="toolbar">
      <el-select v-model="filters.type" clearable placeholder="类型" style="width: 140px">
        <el-option v-for="item in PART_TYPES" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-select v-model="filters.categoryId" clearable placeholder="分类" style="width: 180px">
        <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="load">查询</el-button>
    </div>

    <div class="table-shell">
      <el-table v-loading="loading" :data="items" row-key="id">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="配件" min-width="320">
          <template #default="{ row }">
            <div class="text-strong">{{ row.title }}</div>
            <div class="content-cell text-muted">{{ row.description }}</div>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="90">
          <template #default="{ row }">{{ labelOf(PART_TYPES, row.type) }}</template>
        </el-table-column>
        <el-table-column label="价格" width="110">
          <template #default="{ row }">{{ priceText(row.price) }}</template>
        </el-table-column>
        <el-table-column label="城市" width="120">
          <template #default="{ row }">{{ row.cityName || "-" }}</template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ labelOf(CONTENT_STATUSES, row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="数据" width="150">
          <template #default="{ row }">{{ row.viewCount }} 浏览 / {{ row.intentCount }} 意向</template>
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
import {
  approvePart,
  hidePart,
  listCategories,
  listPartsForReview,
  listPendingReviewParts,
  restorePart
} from "../api/admin";
import type { Part, PartCategory } from "../types";
import { CONTENT_STATUSES, formatDate, labelOf, PART_TYPES, priceText, statusTagType } from "../utils";

const loading = ref(false);
const pendingOnly = ref(false);
const items = ref<Part[]>([]);
const categories = ref<PartCategory[]>([]);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);
const filters = reactive<{ type?: number; categoryId?: number }>({});

async function load() {
  loading.value = true;
  try {
    const params = { page: page.value, pageSize: pageSize.value, type: filters.type, categoryId: filters.categoryId };
    const result = pendingOnly.value ? await listPendingReviewParts(params) : await listPartsForReview(params);
    items.value = result.items;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

async function audit(id: number, action: "approve" | "hide" | "restore") {
  const title = action === "approve" ? "审核通过" : action === "hide" ? "隐藏二手件" : "恢复二手件";
  const { value } = await ElMessageBox.prompt("请输入审核备注", title, {
    inputPlaceholder: "备注",
    confirmButtonText: "确认",
    cancelButtonText: "取消"
  });
  if (action === "approve") {
    await approvePart(id, value || "");
  } else if (action === "hide") {
    await hidePart(id, value || "");
  } else {
    await restorePart(id, value || "");
  }
  await load();
}

onMounted(async () => {
  categories.value = await listCategories();
  await load();
});
</script>
