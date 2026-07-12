<template>
  <main class="login-page">
    <section class="login-panel">
      <div class="login-copy">
        <div class="brand-line">撸车日记 retrofit</div>
        <h1>管理后台</h1>
        <p>用户、内容审核、二手件、举报、车型库和配件分类集中处理。</p>
      </div>
      <form class="login-form" @submit.prevent="submit">
        <label class="field">
          <span>管理员账号</span>
          <input v-model="form.username" class="native-input" autocomplete="username" placeholder="admin" />
        </label>
        <label class="field">
          <span>密码</span>
          <input
            v-model="form.password"
            class="native-input"
            type="password"
            autocomplete="current-password"
            placeholder="开发环境暂未校验密码"
          />
        </label>
        <button class="login-button" type="submit" :disabled="loading">{{ loading ? "登录中..." : "登录" }}</button>
      </form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useAuthStore } from "../stores/auth";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const loading = ref(false);

const form = reactive({
  username: "admin",
  password: "admin"
});

async function submit() {
  if (!form.username.trim()) {
    ElMessage.warning("请输入管理员账号");
    return;
  }
  loading.value = true;
  try {
    await auth.login(form.username, form.password);
    router.replace(String(route.query.redirect || "/"));
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  display: grid;
  min-height: 100vh;
  padding: 24px;
  place-items: center;
  background: #e9eef2;
}

.login-panel {
  display: grid;
  grid-template-columns: 1fr 420px;
  width: min(960px, calc(100vw - 48px));
  min-height: 500px;
  overflow: hidden;
  border: 1px solid #d0d5dd;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 20px 60px rgba(16, 24, 32, 0.12);
}

.login-copy {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 56px;
  background: #101820;
  color: #ffffff;
}

.brand-line {
  color: #79d0c1;
  font-size: 14px;
  font-weight: 700;
}

.login-copy h1 {
  margin: 28px 0 12px;
  font-size: 38px;
  line-height: 1.15;
}

.login-copy p {
  width: 360px;
  margin: 0;
  color: #cbd5e1;
  line-height: 1.75;
}

.login-form {
  align-self: center;
  padding: 48px 40px;
}

.field {
  display: block;
  margin-bottom: 20px;
}

.field span {
  display: block;
  margin-bottom: 8px;
  color: #344054;
  font-size: 14px;
  font-weight: 650;
}

.native-input {
  display: block;
  width: 100%;
  height: 44px;
  padding: 0 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background: #ffffff;
  color: #101820;
  font: inherit;
  outline: none;
}

.native-input:focus {
  border-color: #0f7b6c;
  box-shadow: 0 0 0 3px rgba(15, 123, 108, 0.14);
}

.login-button {
  width: 100%;
  height: 44px;
  border: 0;
  border-radius: 8px;
  background: #0f7b6c;
  color: #ffffff;
  cursor: pointer;
  font: inherit;
  font-weight: 700;
}

.login-button:disabled {
  cursor: not-allowed;
  opacity: 0.72;
}

@media (max-width: 820px) {
  .login-page {
    align-items: stretch;
    padding: 16px;
  }

  .login-panel {
    grid-template-columns: 1fr;
    width: 100%;
    min-height: auto;
  }

  .login-copy {
    padding: 32px 28px;
  }

  .login-copy h1 {
    margin-top: 16px;
    font-size: 30px;
  }

  .login-copy p {
    width: auto;
  }

  .login-form {
    padding: 28px;
  }
}
</style>
