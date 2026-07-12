<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">配件分类</h1>
        <div class="page-desc">维护二手件发布和筛选所用分类。</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">新增分类</el-button>
    </div>

    <div class="table-shell">
      <el-table v-loading="loading" :data="items" row-key="id">
        <el-table-column prop="id" label="ID" width="90" />
        <el-table-column prop="name" label="分类名称" min-width="220" />
        <el-table-column prop="parentId" label="父级 ID" width="120" />
        <el-table-column prop="sortOrder" label="排序" width="120" />
        <el-table-column label="更新时间" width="170">
          <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button size="small" plain @click="openEdit(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑分类' : '新增分类'" width="480px">
      <el-form label-position="top" :model="form">
        <el-form-item label="分类名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="父级 ID">
          <el-input-number v-model="form.parentId" :min="1" controls-position="right" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sortOrder" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import { Plus } from "@element-plus/icons-vue";
import { createCategory, listCategories, updateCategory } from "../api/admin";
import type { PartCategory } from "../types";
import { formatDate } from "../utils";

const loading = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const items = ref<PartCategory[]>([]);
const form = reactive<{ parentId?: number; name: string; sortOrder: number }>({
  parentId: undefined,
  name: "",
  sortOrder: 0
});

async function load() {
  loading.value = true;
  try {
    items.value = await listCategories();
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  editingId.value = null;
  form.parentId = undefined;
  form.name = "";
  form.sortOrder = 0;
  dialogVisible.value = true;
}

function openEdit(row: PartCategory) {
  editingId.value = row.id;
  form.parentId = row.parentId;
  form.name = row.name;
  form.sortOrder = row.sortOrder;
  dialogVisible.value = true;
}

async function submit() {
  if (editingId.value) {
    await updateCategory(editingId.value, form);
  } else {
    await createCategory(form);
  }
  ElMessage.success("已保存");
  dialogVisible.value = false;
  await load();
}

onMounted(load);
</script>
