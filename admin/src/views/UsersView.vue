<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">用户管理</h1>
        <div class="page-desc">查询用户、查看手机号绑定情况，并执行封禁或解封。</div>
      </div>
    </div>

    <div class="toolbar">
      <el-input v-model="filters.q" clearable placeholder="昵称 / 手机号 / OpenID" style="width: 280px" @keyup.enter="load" />
      <el-select v-model="filters.status" clearable placeholder="状态" style="width: 140px">
        <el-option v-for="item in USER_STATUSES" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="load">查询</el-button>
    </div>

    <div class="table-shell">
      <el-table v-loading="loading" :data="items" row-key="id">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户" min-width="180">
          <template #default="{ row }">
            <div class="text-strong">{{ row.nickname || "未命名车友" }}</div>
            <div class="text-muted mono">{{ row.openid }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="手机号" width="150" />
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">{{ labelOf(USER_STATUSES, row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="注册时间" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 1" size="small" type="danger" plain @click="changeStatus(row.id, 2)">封禁</el-button>
            <el-button v-else size="small" type="success" plain @click="changeStatus(row.id, 1)">解封</el-button>
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
import { listUsers, updateUserStatus } from "../api/admin";
import type { User } from "../types";
import { formatDate, labelOf, USER_STATUSES } from "../utils";

const loading = ref(false);
const items = ref<User[]>([]);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);
const filters = reactive<{ q: string; status?: number }>({ q: "" });

async function load() {
  loading.value = true;
  try {
    const result = await listUsers({
      page: page.value,
      pageSize: pageSize.value,
      q: filters.q,
      status: filters.status
    });
    items.value = result.items;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

async function changeStatus(userId: number, status: number) {
  const label = status === 2 ? "封禁" : "解封";
  const { value } = await ElMessageBox.prompt(`请输入${label}备注`, label, {
    inputPlaceholder: "备注",
    confirmButtonText: "确认",
    cancelButtonText: "取消"
  });
  await updateUserStatus(userId, { status, remark: value || "" });
  await load();
}

onMounted(load);
</script>
