<template>
  <div class="dashboard-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">运营概览</h1>
        <div class="page-desc">基于当前后审接口统计待处理工作量。</div>
      </div>
      <el-button :disabled="loading" @click="load">
        <span v-if="loading" class="button-spinner" aria-hidden="true"></span>
        <Refresh v-else class="button-icon" aria-hidden="true" />
        <span>{{ loading ? "刷新中" : "刷新" }}</span>
      </el-button>
    </div>

    <el-alert
      v-if="loadError"
      class="dashboard-alert"
      type="warning"
      :closable="false"
      title="API 暂不可用，已保留页面操作入口。请确认后端服务已启动后刷新统计。"
    />

    <div class="stats-grid">
      <section class="stat-card">
        <div class="stat-label">待审核帖子</div>
        <div class="stat-value">{{ loading ? "-" : stats.pendingPosts }}</div>
      </section>
      <section class="stat-card">
        <div class="stat-label">待审核二手件</div>
        <div class="stat-value">{{ loading ? "-" : stats.pendingParts }}</div>
      </section>
      <section class="stat-card">
        <div class="stat-label">待处理举报</div>
        <div class="stat-value">{{ loading ? "-" : stats.pendingReports }}</div>
      </section>
      <section class="stat-card">
        <div class="stat-label">用户总量</div>
        <div class="stat-value">{{ loading ? "-" : stats.users }}</div>
      </section>
    </div>

    <div class="quick-grid">
      <button class="quick-card primary" type="button" @click="go('/posts')">
        <span>帖子后审</span>
        <strong>处理帖子</strong>
        <em>待审、隐藏、恢复</em>
      </button>
      <button class="quick-card primary" type="button" @click="go('/parts')">
        <span>二手件后审</span>
        <strong>处理二手件</strong>
        <em>出售、求购内容巡检</em>
      </button>
      <button class="quick-card warning" type="button" @click="go('/reports')">
        <span>举报处理</span>
        <strong>处理举报</strong>
        <em>处理或忽略用户举报</em>
      </button>
      <button class="quick-card" type="button" @click="go('/vehicles')">
        <span>车型库</span>
        <strong>维护车型</strong>
        <em>品牌、车系、车型</em>
      </button>
      <button class="quick-card" type="button" @click="go('/categories')">
        <span>配件分类</span>
        <strong>维护分类</strong>
        <em>二手件发布分类</em>
      </button>
      <button class="quick-card" type="button" @click="go('/users')">
        <span>用户管理</span>
        <strong>用户状态</strong>
        <em>查询、封禁、解封</em>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { Refresh } from "@element-plus/icons-vue";
import { listPendingReviewParts, listPendingReviewPosts, listReports, listUsers } from "../api/admin";

const router = useRouter();
const loading = ref(false);
const loadError = ref(false);
const stats = reactive({
  pendingPosts: 0,
  pendingParts: 0,
  pendingReports: 0,
  users: 0
});

async function load() {
  loading.value = true;
  loadError.value = false;
  try {
    const [posts, parts, reports, users] = await Promise.all([
      listPendingReviewPosts({ page: 1, pageSize: 1 }),
      listPendingReviewParts({ page: 1, pageSize: 1 }),
      listReports({ page: 1, pageSize: 1, status: 0 }),
      listUsers({ page: 1, pageSize: 1 })
    ]);
    stats.pendingPosts = posts.total;
    stats.pendingParts = parts.total;
    stats.pendingReports = reports.total;
    stats.users = users.total;
  } catch {
    loadError.value = true;
  } finally {
    loading.value = false;
  }
}

function go(path: string) {
  router.push(path);
}

onMounted(load);
</script>

<style scoped>
.dashboard-page {
  min-width: 0;
}

.dashboard-alert {
  margin-bottom: 16px;
}

.page-header :deep(.el-button > span) {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.button-icon {
  display: block;
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  fill: currentColor;
}

.button-spinner {
  display: block;
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  border: 2px solid #d0d5dd;
  border-top-color: #0f7b6c;
  border-radius: 999px;
  animation: dashboard-spin 0.8s linear infinite;
}

@keyframes dashboard-spin {
  to {
    transform: rotate(360deg);
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(160px, 1fr));
  gap: 16px;
}

.stat-card {
  min-height: 118px;
  padding: 20px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
}

.stat-label {
  color: #667085;
  font-size: 13px;
}

.stat-value {
  margin-top: 12px;
  color: #101820;
  font-size: 34px;
  font-weight: 760;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 18px;
}

.quick-card {
  display: flex;
  min-height: 132px;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  padding: 18px 20px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
  color: #101820;
  cursor: pointer;
  text-align: left;
  transition:
    border-color 0.16s ease,
    box-shadow 0.16s ease,
    transform 0.16s ease;
}

.quick-card:hover {
  border-color: #0f7b6c;
  box-shadow: 0 12px 28px rgba(16, 24, 32, 0.1);
  transform: translateY(-1px);
}

.quick-card span {
  color: #667085;
  font-size: 13px;
}

.quick-card strong {
  margin-top: 10px;
  font-size: 18px;
  font-weight: 760;
}

.quick-card em {
  margin-top: 8px;
  color: #667085;
  font-size: 13px;
  font-style: normal;
}

.quick-card.primary strong {
  color: #0f7b6c;
}

.quick-card.warning strong {
  color: #b54708;
}

@media (max-width: 1100px) {
  .stats-grid,
  .quick-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 680px) {
  .stats-grid,
  .quick-grid {
    grid-template-columns: 1fr;
  }
}
</style>
