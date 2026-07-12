import { defineStore } from "pinia";
import { adminLogin } from "../api/admin";
import type { AdminUser } from "../types";

interface AuthState {
  admin: AdminUser | null;
  token: string;
}

function readAdmin(): AdminUser | null {
  const raw = localStorage.getItem("admin");
  if (!raw) {
    return null;
  }
  try {
    return JSON.parse(raw) as AdminUser;
  } catch {
    return null;
  }
}

export const useAuthStore = defineStore("auth", {
  state: (): AuthState => ({
    admin: readAdmin(),
    token: localStorage.getItem("adminToken") || ""
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.token && state.admin)
  },
  actions: {
    async login(username: string, password: string) {
      const result = await adminLogin({ username, password });
      this.admin = result.admin;
      this.token = result.token;
      localStorage.setItem("admin", JSON.stringify(result.admin));
      localStorage.setItem("adminToken", result.token);
      localStorage.setItem("adminId", String(result.admin.id));
    },
    logout() {
      this.admin = null;
      this.token = "";
      localStorage.removeItem("admin");
      localStorage.removeItem("adminToken");
      localStorage.removeItem("adminId");
    }
  }
});
