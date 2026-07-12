export const API_BASE_URL = "http://127.0.0.1:8080/api/v1";

export function getApiBaseUrl(): string {
  const app = getApp<IAppOption>();
  return app.globalData.apiBaseUrl || API_BASE_URL;
}
