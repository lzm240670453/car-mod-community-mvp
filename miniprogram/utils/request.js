"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ApiError = void 0;
exports.getCurrentUserId = getCurrentUserId;
exports.requireCurrentUserId = requireCurrentUserId;
exports.setCurrentUserId = setCurrentUserId;
exports.request = request;
exports.toQuery = toQuery;
exports.toastError = toastError;
const config_1 = require("./config");
class ApiError extends Error {
    constructor(message, statusCode) {
        super(message);
        this.statusCode = statusCode;
    }
}
exports.ApiError = ApiError;
function getCurrentUserId() {
    const raw = wx.getStorageSync("userId");
    const id = Number(raw);
    return Number.isFinite(id) && id > 0 ? id : 0;
}
function requireCurrentUserId() {
    const userId = getCurrentUserId();
    if (!userId) {
        throw new ApiError("请先在“我的”页面登录", 401);
    }
    return userId;
}
function setCurrentUserId(userId) {
    wx.setStorageSync("userId", userId);
    const app = getApp();
    app.globalData.userId = userId;
}
function request(options) {
    const method = options.method || "GET";
    const url = `${(0, config_1.getApiBaseUrl)()}${options.path}`;
    const userId = getCurrentUserId();
    const header = {
        "Content-Type": "application/json",
        ...options.header
    };
    if (userId) {
        header["X-User-ID"] = String(userId);
    }
    return new Promise((resolve, reject) => {
        wx.request({
            url,
            method,
            data: options.data,
            header,
            success(res) {
                if (res.statusCode === 204) {
                    resolve(undefined);
                    return;
                }
                const body = res.data;
                if (res.statusCode >= 200 && res.statusCode < 300 && body && "data" in body) {
                    resolve(body.data);
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
function toQuery(params) {
    const query = Object.entries(params)
        .filter(([, value]) => value !== undefined && value !== null && value !== "")
        .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
        .join("&");
    return query ? `?${query}` : "";
}
function toastError(error) {
    const message = error instanceof Error ? error.message : "操作失败";
    wx.showToast({ title: message, icon: "none" });
}
