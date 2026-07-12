<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">举报处理</h1>
        <div class="page-desc">按目标类型和处理状态筛选举报，并记录处理备注。</div>
      </div>
    </div>

    <div class="toolbar">
      <el-select v-model="filters.status" clearable placeholder="状态" style="width: 140px">
        <el-option v-for="item in REPORT_STATUSES" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-select v-model="filters.targetType" clearable placeholder="目标类型" style="width: 140px">
        <el-option v-for="item in TARGET_TYPES" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="load">查询</el-button>
    </div>

    <div class="table-shell">
      <el-table v-loading="loading" :data="items" row-key="id">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="reporterId" label="举报人" width="100" />
        <el-table-column label="目标" width="150">
          <template #default="{ row }">{{ labelOf(TARGET_TYPES, row.targetType) }} #{{ row.targetId }}</template>
        </el-table-column>
        <el-table-column prop="reasonText" label="原因" min-width="300" />
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'warning' : 'info'">{{ labelOf(REPORT_STATUSES, row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="提交时间" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="success" plain @click="process(row.id, 1)">处理</el-button>
            <el-button size="small" plain @click="process(row.id, 2)">忽略</el-button>
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
import { listReports, processReport } from "../api/admin";
import type { Report } from "../types";
import { formatDate, labelOf, REPORT_STATUSES, TARGET_TYPES } from "../utils";

const loading = ref(false);
const items = ref<Report[]>([]);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);
const filters = reactive<{ status?: number; targetType?: number }>({ status: 0 });

async function load() {
  loading.value = true;
  try {
    const result = await listReports({
      page: page.value,
      pageSize: pageSize.value,
      status: filters.status,
      targetType: filters.targetType
    });
    items.value = result.items;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

async function process(id: number, status: number) {
  const title = status === 1 ? "处理举报" : "忽略举报";
  const { value } = await ElMessageBox.prompt("请输入处理备注", title, {
    inputPlaceholder: "备注",
    confirmButtonText: "确认",
    cancelButtonText: "取消"
  });
  await processReport(id, { status, remark: value || "" });
  await load();
}

onMounted(load);
</script>
