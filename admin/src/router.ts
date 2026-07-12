import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "./stores/auth";
import LoginView from "./views/LoginView.vue";
import DashboardView from "./views/DashboardView.vue";
import UsersView from "./views/UsersView.vue";
import PostsReviewView from "./views/PostsReviewView.vue";
import PartsReviewView from "./views/PartsReviewView.vue";
import ReportsView from "./views/ReportsView.vue";
import VehiclesView from "./views/VehiclesView.vue";
import CategoriesView from "./views/CategoriesView.vue";

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/login", name: "login", component: LoginView, meta: { public: true } },
    { path: "/", name: "dashboard", component: DashboardView },
    { path: "/users", name: "users", component: UsersView },
    { path: "/posts", name: "posts", component: PostsReviewView },
    { path: "/parts", name: "parts", component: PartsReviewView },
    { path: "/reports", name: "reports", component: ReportsView },
    { path: "/vehicles", name: "vehicles", component: VehiclesView },
    { path: "/categories", name: "categories", component: CategoriesView }
  ]
});

router.beforeEach((to) => {
  const auth = useAuthStore();
  if (!to.meta.public && !auth.isLoggedIn) {
    return { path: "/login", query: { redirect: to.fullPath } };
  }
  if (to.path === "/login" && auth.isLoggedIn) {
    return "/";
  }
  return true;
});
