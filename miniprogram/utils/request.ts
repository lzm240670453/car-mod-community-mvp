import { ApiResponse } from "../types/index";
import { getApiBaseUrl } from "./config";

type Method = "GET" | "POST" | "PUT" | "DELETE";

interface RequestOptions {
  path: string;
  method?: Method;
  data?: WechatMiniprogram.IAnyObject | string | number | boolean | null;
  header?: Record<string, string>;
}

export class ApiError extends Error {
  statusCode: number;

  constructor(message: string, statusCode: number) {
    super(message);
    this.statusCode = statusCode;
  }
}

export function getCurrentUserId(): number {
  const raw = wx.getStorageSync("userId");
  const id = Number(raw);
  return Number.isFinite(id) && id > 0 ? id : 0;
}

export function requireCurrentUserId(): number {
  const userId = getCurrentUserId();
  if (!userId) {
    throw new ApiError("请先在“我的”页面登录", 401);
  }
  return userId;
}

export function setCurrentUserId(userId: number) {
  wx.setStorageSync("userId", userId);
  const app = getApp<IAppOption>();
  app.globalData.userId = userId;
}

export function request<T>(options: RequestOptions): Promise<T> {
  const method = options.method || "GET";
  const url = `${getApiBaseUrl()}${options.path}`;
  const userId = getCurrentUserId();
  const header: Record<string, string> = {
    "Content-Type": "application/json",
    ...options.header
  };

  if (userId) {
    header["X-User-ID"] = String(userId);
  }

  return new Promise<T>((resolve, reject) => {
    wx.request<ApiResponse<T>>({
      url,
      method,
      data: options.data as WechatMiniprogram.IAnyObject,
      header,
      success(res) {
        if (res.statusCode === 204) {
          resolve(undefined as T);
          return;
        }

        const body = res.data;
        if (res.statusCode >= 200 && res.statusCode < 300 && body && "data" in body) {
          resolve(body.data as T);
          return;
        }

        const message = body?.error || `请求失败：${res.statusCode}`;
        reject(new ApiError(message, res.statusCode));
      },
      fail(err) {
        reject(new ApiError(err.errMsg || "网络请求失败", 0));
      }
    });
  });
}

export function toQuery(params: Record<string, string | number | undefined | null>): string {
  const query = Object.entries(params)
    .filter(([, value]) => value !== undefined && value !== null && value !== "")
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
    .join("&");
  return query ? `?${query}` : "";
}

export function toastError(error: unknown) {
  const message = error instanceof Error ? error.message : "操作失败";
  wx.showToast({ title: message, icon: "none" });
}
