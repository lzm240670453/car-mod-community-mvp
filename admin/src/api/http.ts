import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";
import { ElMessage } from "element-plus";
import type { ApiResponse } from "../types";

export const http = axios.create({
  baseURL: "/api/v1",
  timeout: 15000
});

http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const adminId = localStorage.getItem("adminId") || "1";
  config.headers.set("X-Admin-ID", adminId);
  return config;
});

http.interceptors.response.use(
  (response) => {
    if (response.status === 204) {
      return undefined;
    }
    const body = response.data as ApiResponse<unknown>;
    if (body && "data" in body) {
      return body.data;
    }
    return response.data;
  },
  (error: AxiosError<ApiResponse<unknown>>) => {
    const message = error.response?.data?.error || error.message || "请求失败";
    ElMessage.error(message);
    return Promise.reject(new Error(message));
  }
);
