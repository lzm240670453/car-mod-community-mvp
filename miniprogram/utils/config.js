"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.API_BASE_URL = void 0;
exports.getApiBaseUrl = getApiBaseUrl;
exports.API_BASE_URL = "http://127.0.0.1:8080/api/v1";
function getApiBaseUrl() {
    const app = getApp();
    return app.globalData.apiBaseUrl || exports.API_BASE_URL;
}
