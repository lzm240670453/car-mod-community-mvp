export interface ApiResponse<T> {
  data?: T;
  error?: string;
}

export interface PageResult<T> {
  items: T[];
  page: number;
  pageSize: number;
  total: number;
}

export interface Timestamp {
  createdAt: string;
  updatedAt: string;
}

export interface AdminUser extends Timestamp {
  id: number;
  username: string;
  role: string;
  status: number;
}

export interface User extends Timestamp {
  id: number;
  openid: string;
  unionid?: string;
  nickname: string;
  avatarUrl: string;
  phone: string;
  phoneBoundAt?: string;
  status: number;
}

export interface Post extends Timestamp {
  id: number;
  userId: number;
  garageId?: number;
  type: number;
  title: string;
  content: string;
  vehicleModelId?: number;
  status: number;
  likeCount: number;
  commentCount: number;
  favoriteCount: number;
}

export interface Part extends Timestamp {
  id: number;
  userId: number;
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
  status: number;
  viewCount: number;
  favoriteCount: number;
  intentCount: number;
}

export interface Report extends Timestamp {
  id: number;
  reporterId: number;
  targetType: number;
  targetId: number;
  reasonType: number;
  reasonText: string;
  status: number;
}

export interface VehicleAdminItem {
  modelId: number;
  modelName: string;
  year?: number;
  generation: string;
  engine: string;
  seriesId: number;
  seriesName: string;
  brandId: number;
  brandName: string;
}

export interface VehicleBrand extends Timestamp {
  id: number;
  name: string;
  initial: string;
  logoUrl: string;
  sortOrder: number;
}

export interface VehicleSeries extends Timestamp {
  id: number;
  brandId: number;
  name: string;
}

export interface VehicleModel extends Timestamp {
  id: number;
  seriesId: number;
  name: string;
  year?: number;
  generation: string;
  engine: string;
}

export interface PartCategory extends Timestamp {
  id: number;
  parentId?: number;
  name: string;
  sortOrder: number;
}
