import {
  Comment,
  PageResult,
  Part,
  PartCategory,
  PartDetail,
  Post,
  PostDetail,
  Report,
  SiteMessage,
  TradeIntent,
  User,
  UserGarage,
  VehicleBrand,
  VehicleModel,
  VehicleSearchItem,
  VehicleSeries
} from "../types/index";
import { request, toQuery } from "../utils/request";

export interface WechatLoginPayload {
  code?: string;
  openid?: string;
  unionid?: string;
  nickname?: string;
  avatarUrl?: string;
}

export function wechatLogin(payload: WechatLoginPayload) {
  return request<{ user: User; token: string }>({
    path: "/auth/wechat-login",
    method: "POST",
    data: payload
  });
}

export function getMe() {
  return request<User>({ path: "/users/me" });
}

export function updateMe(payload: { nickname: string; avatarUrl: string }) {
  return request<User>({ path: "/users/me", method: "PUT", data: payload });
}

export function bindPhone(phone: string) {
  return request<User>({ path: "/auth/bind-phone", method: "POST", data: { phone } });
}

export function listGarages() {
  return request<UserGarage[]>({ path: "/users/me/garages" });
}

export function createGarage(payload: {
  vehicleModelId: number;
  year?: number;
  nickname: string;
  description: string;
  isPrimary: boolean;
}) {
  return request<UserGarage>({ path: "/users/me/garages", method: "POST", data: payload });
}

export function listVehicleBrands() {
  return request<VehicleBrand[]>({ path: "/vehicles/brands" });
}

export function listVehicleSeries(brandId: number) {
  return request<VehicleSeries[]>({ path: `/vehicles/brands/${brandId}/series` });
}

export function listVehicleModels(seriesId: number) {
  return request<VehicleModel[]>({ path: `/vehicles/series/${seriesId}/models` });
}

export function searchVehicles(q: string) {
  return request<VehicleSearchItem[]>({ path: `/vehicles/search${toQuery({ q })}` });
}

export function listPosts(params: {
  page?: number;
  pageSize?: number;
  type?: number;
  q?: string;
  userId?: number;
  vehicleModelId?: number;
}) {
  return request<PageResult<Post>>({ path: `/posts${toQuery(params)}` });
}

export function getPost(postId: number) {
  return request<PostDetail>({ path: `/posts/${postId}` });
}

export function createPost(payload: {
  garageId?: number;
  type: number;
  title: string;
  content: string;
  vehicleModelId?: number;
  images: string[];
}) {
  return request<{ id: number; status: number }>({ path: "/posts", method: "POST", data: payload });
}

export function likePost(postId: number) {
  return request<{ id: number }>({ path: `/posts/${postId}/like`, method: "POST" });
}

export function favoritePost(postId: number) {
  return request<{ id: number }>({ path: `/posts/${postId}/favorite`, method: "POST" });
}

export function listPostComments(postId: number, params: { page?: number; pageSize?: number }) {
  return request<PageResult<Comment>>({ path: `/posts/${postId}/comments${toQuery(params)}` });
}

export function createComment(postId: number, payload: { parentId?: number; content: string }) {
  return request<{ id: number; status: number }>({
    path: `/posts/${postId}/comments`,
    method: "POST",
    data: payload
  });
}

export function listPartCategories() {
  return request<PartCategory[]>({ path: "/part-categories" });
}

export function listParts(params: {
  page?: number;
  pageSize?: number;
  type?: number;
  categoryId?: number;
  cityCode?: string;
  conditionLevel?: number;
  minPrice?: number;
  maxPrice?: number;
  q?: string;
  userId?: number;
  vehicleModelId?: number;
}) {
  return request<PageResult<Part>>({ path: `/parts${toQuery(params)}` });
}

export function getPart(partId: number) {
  return request<PartDetail>({ path: `/parts/${partId}` });
}

export function createPart(payload: {
  type: number;
  categoryId: number;
  title: string;
  brand: string;
  model: string;
  conditionLevel: number;
  price?: number;
  cityCode: string;
  cityName: string;
  description: string;
  contactPolicy: number;
  images: string[];
  fitments: Array<{ vehicleModelId: number; note: string }>;
}) {
  return request<{ id: number; status: number }>({ path: "/parts", method: "POST", data: payload });
}

export function favoritePart(partId: number) {
  return request<{ id: number }>({ path: `/parts/${partId}/favorite`, method: "POST" });
}

export function createIntent(partId: number, message: string) {
  return request<{ intentId: number; status: number }>({
    path: `/parts/${partId}/intents`,
    method: "POST",
    data: { message }
  });
}

export function listIntents(params: { page?: number; pageSize?: number; status?: number; role?: string }) {
  return request<PageResult<TradeIntent>>({ path: `/intents${toQuery(params)}` });
}

export function closeIntent(intentId: number) {
  return request<void>({ path: `/intents/${intentId}/close`, method: "POST" });
}

export function listMessages(params: { page?: number; pageSize?: number; type?: number; unread?: number }) {
  return request<PageResult<SiteMessage>>({ path: `/messages${toQuery(params)}` });
}

export function getUnreadMessageCount() {
  return request<{ count: number }>({ path: "/messages/unread-count" });
}

export function markMessageRead(messageId: number) {
  return request<void>({ path: `/messages/read/${messageId}`, method: "POST" });
}

export function markAllMessagesRead() {
  return request<void>({ path: "/messages/read-all", method: "POST" });
}

export function createReport(payload: {
  targetType: number;
  targetId: number;
  reasonType: number;
  reasonText: string;
}) {
  return request<{ id: number; status: number }>({ path: "/reports", method: "POST", data: payload });
}

export function listMyReports(params: { page?: number; pageSize?: number }) {
  return request<PageResult<Report>>({ path: `/reports/my${toQuery(params)}` });
}
