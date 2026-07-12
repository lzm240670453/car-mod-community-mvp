import { http } from "./http";
import type {
  AdminUser,
  PageResult,
  Part,
  PartCategory,
  Post,
  Report,
  User,
  VehicleAdminItem,
  VehicleBrand,
  VehicleModel,
  VehicleSeries
} from "../types";

export function adminLogin(payload: { username: string; password: string }) {
  return http.post<unknown, { admin: AdminUser; token: string }>("/admin/auth/login", payload);
}

export function listUsers(params: { page?: number; pageSize?: number; status?: number; q?: string }) {
  return http.get<unknown, PageResult<User>>("/admin/users", { params });
}

export function updateUserStatus(userId: number, payload: { status: number; remark: string }) {
  return http.put<unknown, { id: number; status: number }>(`/admin/users/${userId}/status`, payload);
}

export function listPostsForReview(params: { page?: number; pageSize?: number; type?: number; userId?: number }) {
  return http.get<unknown, PageResult<Post>>("/admin/posts/review", { params });
}

export function listPendingReviewPosts(params: { page?: number; pageSize?: number; type?: number; userId?: number }) {
  return http.get<unknown, PageResult<Post>>("/admin/posts/pending-review", { params });
}

export function approvePost(postId: number, remark = "") {
  return http.post<unknown, { id: number; status: number }>(`/admin/posts/${postId}/approve`, { remark });
}

export function hidePost(postId: number, remark = "") {
  return http.post<unknown, { id: number; status: number }>(`/admin/posts/${postId}/hide`, { remark });
}

export function restorePost(postId: number, remark = "") {
  return http.post<unknown, { id: number; status: number }>(`/admin/posts/${postId}/restore`, { remark });
}

export function listPartsForReview(params: { page?: number; pageSize?: number; type?: number; categoryId?: number }) {
  return http.get<unknown, PageResult<Part>>("/admin/parts/review", { params });
}

export function listPendingReviewParts(params: { page?: number; pageSize?: number; type?: number; categoryId?: number }) {
  return http.get<unknown, PageResult<Part>>("/admin/parts/pending-review", { params });
}

export function approvePart(partId: number, remark = "") {
  return http.post<unknown, { id: number; status: number }>(`/admin/parts/${partId}/approve`, { remark });
}

export function hidePart(partId: number, remark = "") {
  return http.post<unknown, { id: number; status: number }>(`/admin/parts/${partId}/hide`, { remark });
}

export function restorePart(partId: number, remark = "") {
  return http.post<unknown, { id: number; status: number }>(`/admin/parts/${partId}/restore`, { remark });
}

export function listReports(params: { page?: number; pageSize?: number; status?: number; targetType?: number }) {
  return http.get<unknown, PageResult<Report>>("/admin/reports", { params });
}

export function processReport(reportId: number, payload: { status: number; remark: string }) {
  return http.post<unknown, { id: number; status: number }>(`/admin/reports/${reportId}/process`, payload);
}

export function listVehicles(params: { page?: number; pageSize?: number; q?: string }) {
  return http.get<unknown, PageResult<VehicleAdminItem>>("/admin/vehicles", { params });
}

export function createBrand(payload: { name: string; initial: string; logoUrl: string; sortOrder: number }) {
  return http.post<unknown, VehicleBrand>("/admin/vehicles/brands", payload);
}

export function createSeries(payload: { brandId: number; name: string }) {
  return http.post<unknown, VehicleSeries>("/admin/vehicles/series", payload);
}

export function createModel(payload: {
  seriesId: number;
  name: string;
  year?: number;
  generation: string;
  engine: string;
}) {
  return http.post<unknown, VehicleModel>("/admin/vehicles/models", payload);
}

export function listCategories() {
  return http.get<unknown, PartCategory[]>("/admin/categories");
}

export function createCategory(payload: { parentId?: number; name: string; sortOrder: number }) {
  return http.post<unknown, PartCategory>("/admin/categories", payload);
}

export function updateCategory(categoryId: number, payload: { parentId?: number; name: string; sortOrder: number }) {
  return http.put<unknown, PartCategory>(`/admin/categories/${categoryId}`, payload);
}
