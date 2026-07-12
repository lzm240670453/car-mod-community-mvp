<template>
  <RouterView v-if="route.meta.public" />
  <el-container v-else class="app-shell">
    <el-aside class="app-aside">
      <div class="brand">
        <div class="brand-mark">R</div>
        <div>
          <div class="brand-title">撸车日记</div>
          <div class="brand-subtitle">管理后台</div>
        </div>
      </div>
      <nav class="nav-menu" aria-label="后台导航">
        <RouterLink
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ 'is-active': route.path === item.path }"
          :title="item.label"
        >
          <span class="nav-icon" aria-hidden="true">
            <component :is="item.icon" />
          </span>
          <span class="nav-label">{{ item.label }}</span>
        </RouterLink>
      </nav>
    </el-aside>
    <el-container class="app-content">
      <el-header class="app-header">
        <div class="header-left">
          <span class="status-dot ok"></span>
          <span>API /api/v1</span>
        </div>
        <div class="header-right">
          <span class="text-muted">{{ auth.admin?.username }}</span>
          <button class="logout-button" type="button" @click="logout">
            <SwitchButton class="header-button-icon" aria-hidden="true" />
            <span>退出</span>
          </button>
        </div>
      </el-header>
      <el-main class="app-main">
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";
import { DataBoard, Document, Goods, Grid, SwitchButton, User, Van, Warning } from "@element-plus/icons-vue";
import { useAuthStore } from "./stores/auth";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const navItems = [
  { path: "/", label: "概览", icon: DataBoard },
  { path: "/users", label: "用户管理", icon: User },
  { path: "/posts", label: "帖子后审", icon: Document },
  { path: "/parts", label: "二手件后审", icon: Goods },
  { path: "/reports", label: "举报处理", icon: Warning },
  { path: "/vehicles", label: "车型库", icon: Van },
  { path: "/categories", label: "配件分类", icon: Grid }
];

function logout() {
  auth.logout();
  router.push("/login");
}
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  background: #f4f6f8;
}

.app-aside {
  position: sticky;
  top: 0;
  width: 232px;
  min-width: 232px;
  height: 100vh;
  flex-shrink: 0;
  border-right: 1px solid #e5e7eb;
  background: #101820;
  overflow: hidden;
}

.app-content {
  min-width: 0;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 72px;
  padding: 0 20px;
  color: #ffffff;
}

.brand-mark {
  display: grid;
  width: 36px;
  height: 36px;
  place-items: center;
  border-radius: 8px;
  background: #0f7b6c;
  font-weight: 800;
}

.brand-title {
  font-size: 17px;
  font-weight: 750;
}

.brand-subtitle {
  margin-top: 2px;
  color: #a7b0ba;
  font-size: 12px;
}

.nav-menu {
  display: flex;
  width: 100%;
  min-width: 0;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  height: 44px;
  min-width: 0;
  padding: 0 14px;
  border-radius: 8px;
  color: #cbd5e1;
  line-height: 44px;
}

.nav-icon {
  display: inline-flex;
  width: 18px;
  height: 18px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
  color: inherit;
  line-height: 1;
}

.nav-icon :deep(svg) {
  display: block;
  width: 18px;
  height: 18px;
  fill: currentColor;
}

.nav-label {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nav-item:hover,
.nav-item.is-active {
  background: rgba(15, 123, 108, 0.18);
  color: #ffffff;
}

.app-header {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 24px;
  border-bottom: 1px solid #e5e7eb;
  background: #ffffff;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logout-button {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-height: 32px;
  padding: 6px 10px;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: #667085;
  cursor: pointer;
  font: inherit;
  line-height: 1;
}

.logout-button:hover {
  background: #f2f4f7;
  color: #101820;
}

.logout-button:focus-visible {
  outline: 2px solid rgba(15, 123, 108, 0.36);
  outline-offset: 2px;
}

.header-button-icon {
  display: block;
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  fill: currentColor;
}

.app-main {
  min-width: 0;
  padding: 24px;
}

@media (max-width: 900px) {
  .app-aside {
    width: 72px;
    min-width: 72px;
  }

  .brand {
    justify-content: center;
    padding: 0 12px;
  }

  .brand-title,
  .brand-subtitle,
  .nav-label {
    display: none;
  }

  .nav-menu {
    padding: 8px 10px;
  }

  .nav-item {
    justify-content: center;
    width: 52px;
    padding: 0;
  }

  .nav-icon {
    margin-right: 0;
  }

  .app-header {
    padding: 0 14px;
  }

  .header-left {
    display: none;
  }

  .app-main {
    padding: 16px;
  }
}
</style>
