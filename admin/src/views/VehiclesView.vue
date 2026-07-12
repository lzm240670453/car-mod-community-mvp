<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">车型库</h1>
        <div class="page-desc">新增品牌、车系、车型，并查看当前车型库。</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="dialogVisible = true">新增资料</el-button>
    </div>

    <div class="toolbar">
      <el-input v-model="q" clearable placeholder="品牌 / 车系 / 车型" style="width: 280px" @keyup.enter="load" />
      <el-button type="primary" :icon="Search" @click="load">查询</el-button>
    </div>

    <div class="table-shell">
      <el-table v-loading="loading" :data="items" row-key="modelId">
        <el-table-column prop="brandName" label="品牌" width="140" />
        <el-table-column prop="seriesName" label="车系" width="160" />
        <el-table-column prop="modelName" label="车型" min-width="220" />
        <el-table-column prop="year" label="年款" width="100" />
        <el-table-column prop="generation" label="代际" width="140" />
        <el-table-column prop="engine" label="发动机" width="140" />
        <el-table-column label="ID" width="200">
          <template #default="{ row }">
            <span class="mono">B{{ row.brandId }} / S{{ row.seriesId }} / M{{ row.modelId }}</span>
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

    <el-dialog v-model="dialogVisible" title="新增车型资料" width="min(680px, calc(100vw - 32px))">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="品牌" name="brand">
          <el-form label-position="top" :model="brandForm">
            <div class="form-grid">
              <el-form-item label="品牌名称">
                <el-input v-model="brandForm.name" />
              </el-form-item>
              <el-form-item label="首字母">
                <el-input v-model="brandForm.initial" maxlength="1" />
              </el-form-item>
              <el-form-item label="Logo URL">
                <el-input v-model="brandForm.logoUrl" />
              </el-form-item>
              <el-form-item label="排序">
                <el-input-number v-model="brandForm.sortOrder" :min="0" />
              </el-form-item>
            </div>
            <el-button type="primary" @click="submitBrand">保存品牌</el-button>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="车系" name="series">
          <el-form label-position="top" :model="seriesForm">
            <el-form-item label="品牌 ID">
              <el-input-number v-model="seriesForm.brandId" :min="1" />
            </el-form-item>
            <el-form-item label="车系名称">
              <el-input v-model="seriesForm.name" />
            </el-form-item>
            <el-button type="primary" @click="submitSeries">保存车系</el-button>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="车型" name="model">
          <el-form label-position="top" :model="modelForm">
            <div class="form-grid">
              <el-form-item label="车系 ID">
                <el-input-number v-model="modelForm.seriesId" :min="1" />
              </el-form-item>
              <el-form-item label="车型名称">
                <el-input v-model="modelForm.name" />
              </el-form-item>
              <el-form-item label="年款">
                <el-input-number v-model="modelForm.year" :min="1900" :max="2100" />
              </el-form-item>
              <el-form-item label="代际">
                <el-input v-model="modelForm.generation" />
              </el-form-item>
              <el-form-item label="发动机">
                <el-input v-model="modelForm.engine" />
              </el-form-item>
            </div>
            <el-button type="primary" @click="submitModel">保存车型</el-button>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import { Plus, Search } from "@element-plus/icons-vue";
import { createBrand, createModel, createSeries, listVehicles } from "../api/admin";
import type { VehicleAdminItem } from "../types";

const loading = ref(false);
const dialogVisible = ref(false);
const activeTab = ref("brand");
const items = ref<VehicleAdminItem[]>([]);
const q = ref("");
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);

const brandForm = reactive({ name: "", initial: "", logoUrl: "", sortOrder: 0 });
const seriesForm = reactive({ brandId: 1, name: "" });
const modelForm = reactive({ seriesId: 1, name: "", year: undefined as number | undefined, generation: "", engine: "" });

async function load() {
  loading.value = true;
  try {
    const result = await listVehicles({ page: page.value, pageSize: pageSize.value, q: q.value });
    items.value = result.items;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

async function submitBrand() {
  await createBrand(brandForm);
  ElMessage.success("品牌已创建");
  brandForm.name = "";
  brandForm.initial = "";
  brandForm.logoUrl = "";
  brandForm.sortOrder = 0;
  await load();
}

async function submitSeries() {
  await createSeries(seriesForm);
  ElMessage.success("车系已创建");
  seriesForm.name = "";
  await load();
}

async function submitModel() {
  await createModel(modelForm);
  ElMessage.success("车型已创建");
  modelForm.name = "";
  modelForm.generation = "";
  modelForm.engine = "";
  await load();
}

onMounted(load);
</script>
