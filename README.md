# 汽车改装社区 MVP 方案

## 1. 产品定位

产品名称为「撸车记」，英文名为 retrofit。

第一版只做两件事：

- 改装交流：发作业、分享经验、求推荐、避坑讨论。
- 信息撮合：发布/查找二手改装件，平台只撮合线索，不做支付、担保、物流和售后。

目标是让用户低成本玩车，先验证内容和撮合需求，再逐步扩展。

第一版不限制重点车型，不做官方车型推荐；首发城市也不做限制。

## 2. 目标用户

- 普通玩车用户：想看案例、找方案、买便宜可靠的二手改装件。
- 经验用户：愿意分享改装历程、预算、安装过程和避坑经验。
- 改装店/安装店：第一版隐藏，不开放商家体系，只做普通用户。

## 3. MVP 范围

### 3.1 小程序端

- 微信登录
- 用户资料
- 手机号绑定
- 个人车库
- 车型选择
- 社区帖子列表
- 帖子详情
- 发布帖子
- 评论
- 点赞
- 收藏
- 二手件列表
- 二手件详情
- 发布出售
- 发布求购
- 按车型/类目/城市/价格筛选
- 发起交易意向
- 举报内容
- 我的发布
- 我的收藏
- 我的意向
- 帖子和评论中允许出现手机号与微信号

### 3.2 管理后台

- 管理员登录
- 用户管理
- 帖子后审
- 二手件后审
- 举报处理
- 车型库管理
- 配件类目管理
- 敏感词配置
- 审核日志
- 先发后审内容巡检

### 3.3 第一版不做

- 微信支付
- 平台担保交易
- 订单系统
- 退款售后
- 物流跟踪
- 交易评价
- 站内信
- 站内留言
- 商家入驻
- 直播
- 拍卖
- App
- H5 多端
- 复杂推荐算法

## 4. 核心用户流程

### 4.1 发布改装内容

1. 用户微信登录。
2. 用户完善个人车库。
3. 用户选择帖子类型。
4. 用户上传图片、填写标题、内容、车型、配件类目。
5. 内容直接进入信息流。
6. 系统进入后审巡检流程。
7. 其他用户评论、点赞、收藏。

### 4.2 发布二手改装件

1. 用户选择出售或求购。
2. 用户填写配件类目、品牌、型号、适配车型、成色、价格、城市。
3. 用户上传实物图片。
4. 信息直接展示。
5. 系统进入后审巡检流程。
6. 买家点击「我想要」或「联系车友」。
7. 系统生成交易意向。
8. 双方通过帖子、评论中公开的手机号或微信号自行沟通。

### 4.3 举报处理

1. 用户举报帖子、评论或二手件。
2. 后台收到举报。
3. 管理员查看内容和历史记录。
4. 管理员选择通过、下架、封禁、忽略。
5. 系统记录审核日志。

## 5. 页面结构

### 5.1 小程序 Tab

```text
首页
二手件
发布
我的
```

### 5.2 首页

- 推荐帖子流
- 关注车型筛选
- 帖子类型筛选
- 热门车型入口
- 搜索入口

帖子类型：

- 改装作业
- 避坑经验
- 求推荐
- 问题求助
- 车友闲聊

### 5.3 二手件

- 出售列表
- 求购列表
- 车型筛选
- 类目筛选
- 城市筛选
- 价格区间
- 成色筛选
- 只看同城

配件类目第一版：

- 轮毂轮胎
- 避震悬挂
- 刹车
- 排气
- 进气
- 外观套件
- 内饰
- 灯光
- 电子设备
- 其他

### 5.4 发布

发布入口分三类：

- 发改装内容
- 出售改装件
- 求购改装件

### 5.5 我的

- 我的车库
- 我的帖子
- 我的二手件
- 我的求购
- 我的收藏
- 我的意向
- 设置
- 举报记录

## 6. 技术选型

### 6.1 小程序

```text
微信原生小程序
TypeScript
TDesign Miniprogram
npm
```

### 6.2 后端

```text
Go
Gin
GORM
MySQL 8
Redis
JWT + 微信 openid/session
Zap 日志
Swagger/OpenAPI
```

### 6.3 存储

```text
腾讯云 COS
```

上传流程：

1. 小程序请求后端获取上传凭证。
2. 小程序直传 COS。
3. 小程序把文件 key 回传后端。
4. 后端保存图片元数据和业务关联。

### 6.4 后台

```text
Vue 3
Vite
Element Plus
Axios
Pinia
```

### 6.5 部署

MVP 起步：

```text
Docker Compose
Nginx
Go API
MySQL
Redis
```

后期可演进：

```text
对象存储 CDN
独立搜索服务
消息队列
Kubernetes
```

## 7. 后端模块

```text
user      用户、微信登录、资料
garage    用户车库
vehicle   品牌、车系、车型、年款
post      帖子、评论、点赞、收藏
part      二手件、求购、适配关系
intent    交易意向
audit     审核、举报、敏感词
upload    上传凭证和图片元数据
admin     后台管理
```

## 8. 数据库设计

### 8.1 users

```sql
CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  openid VARCHAR(64) NOT NULL UNIQUE,
  unionid VARCHAR(64) NULL,
  nickname VARCHAR(64) NOT NULL DEFAULT '',
  avatar_url VARCHAR(512) NOT NULL DEFAULT '',
  phone VARCHAR(32) NOT NULL DEFAULT '',
  phone_bound_at DATETIME NULL,
  status TINYINT NOT NULL DEFAULT 1 COMMENT '1 normal, 2 banned',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

### 8.2 user_garages

```sql
CREATE TABLE user_garages (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  vehicle_model_id BIGINT NOT NULL,
  year SMALLINT NULL,
  nickname VARCHAR(64) NOT NULL DEFAULT '',
  description VARCHAR(512) NOT NULL DEFAULT '',
  is_primary TINYINT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  INDEX idx_user_id (user_id),
  INDEX idx_vehicle_model_id (vehicle_model_id)
);
```

### 8.3 vehicle_brands

```sql
CREATE TABLE vehicle_brands (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  initial CHAR(1) NOT NULL DEFAULT '',
  logo_url VARCHAR(512) NOT NULL DEFAULT '',
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  UNIQUE KEY uk_name (name)
);
```

### 8.4 vehicle_series

```sql
CREATE TABLE vehicle_series (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  brand_id BIGINT NOT NULL,
  name VARCHAR(64) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  INDEX idx_brand_id (brand_id),
  UNIQUE KEY uk_brand_series (brand_id, name)
);
```

### 8.5 vehicle_models

```sql
CREATE TABLE vehicle_models (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  series_id BIGINT NOT NULL,
  name VARCHAR(128) NOT NULL,
  year SMALLINT NULL,
  generation VARCHAR(64) NOT NULL DEFAULT '',
  engine VARCHAR(64) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  INDEX idx_series_id (series_id)
);
```

### 8.6 posts

```sql
CREATE TABLE posts (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  garage_id BIGINT NULL,
  type TINYINT NOT NULL COMMENT '1 build, 2 pitfall, 3 recommendation, 4 question, 5 chat',
  title VARCHAR(120) NOT NULL,
  content TEXT NOT NULL,
  vehicle_model_id BIGINT NULL,
  status TINYINT NOT NULL DEFAULT 1 COMMENT '0 pending_review, 1 visible, 2 hidden, 3 deleted',
  like_count INT NOT NULL DEFAULT 0,
  comment_count INT NOT NULL DEFAULT 0,
  favorite_count INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  INDEX idx_user_id (user_id),
  INDEX idx_vehicle_model_id (vehicle_model_id),
  INDEX idx_status_created_at (status, created_at)
);
```

### 8.7 post_images

```sql
CREATE TABLE post_images (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL,
  image_url VARCHAR(512) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  INDEX idx_post_id (post_id)
);
```

### 8.8 comments

```sql
CREATE TABLE comments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  parent_id BIGINT NULL,
  content VARCHAR(1000) NOT NULL,
  status TINYINT NOT NULL DEFAULT 1 COMMENT '0 pending_review, 1 visible, 2 hidden, 3 deleted',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  INDEX idx_post_id (post_id),
  INDEX idx_user_id (user_id)
);
```

### 8.9 post_likes

```sql
CREATE TABLE post_likes (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at DATETIME NOT NULL,
  UNIQUE KEY uk_post_user (post_id, user_id),
  INDEX idx_user_id (user_id)
);
```

### 8.10 post_favorites

```sql
CREATE TABLE post_favorites (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at DATETIME NOT NULL,
  UNIQUE KEY uk_post_user (post_id, user_id),
  INDEX idx_user_id (user_id)
);
```

### 8.11 part_categories

```sql
CREATE TABLE part_categories (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  parent_id BIGINT NULL,
  name VARCHAR(64) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

### 8.12 parts

```sql
CREATE TABLE parts (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  type TINYINT NOT NULL COMMENT '1 sell, 2 buy',
  category_id BIGINT NOT NULL,
  title VARCHAR(120) NOT NULL,
  brand VARCHAR(64) NOT NULL DEFAULT '',
  model VARCHAR(128) NOT NULL DEFAULT '',
  condition_level TINYINT NOT NULL DEFAULT 0 COMMENT '0 unknown, 1 new, 2 almost_new, 3 used, 4 repaired',
  price DECIMAL(10, 2) NULL,
  city_code VARCHAR(32) NOT NULL DEFAULT '',
  city_name VARCHAR(64) NOT NULL DEFAULT '',
  description TEXT NOT NULL,
  contact_policy TINYINT NOT NULL DEFAULT 1 COMMENT '1 public_in_content',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '0 pending_review, 1 visible, 2 hidden, 3 deleted, 4 closed',
  view_count INT NOT NULL DEFAULT 0,
  favorite_count INT NOT NULL DEFAULT 0,
  intent_count INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  INDEX idx_user_id (user_id),
  INDEX idx_category_id (category_id),
  INDEX idx_status_created_at (status, created_at),
  INDEX idx_city_code (city_code)
);
```

### 8.13 part_images

```sql
CREATE TABLE part_images (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  part_id BIGINT NOT NULL,
  image_url VARCHAR(512) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  INDEX idx_part_id (part_id)
);
```

### 8.14 part_fitments

```sql
CREATE TABLE part_fitments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  part_id BIGINT NOT NULL,
  vehicle_model_id BIGINT NOT NULL,
  note VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL,
  UNIQUE KEY uk_part_vehicle (part_id, vehicle_model_id),
  INDEX idx_vehicle_model_id (vehicle_model_id)
);
```

### 8.15 part_favorites

```sql
CREATE TABLE part_favorites (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  part_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at DATETIME NOT NULL,
  UNIQUE KEY uk_part_user (part_id, user_id),
  INDEX idx_user_id (user_id)
);
```

### 8.16 trade_intents

```sql
CREATE TABLE trade_intents (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  part_id BIGINT NOT NULL,
  buyer_id BIGINT NOT NULL,
  seller_id BIGINT NOT NULL,
  message VARCHAR(500) NOT NULL DEFAULT '',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '1 active, 2 cancelled, 3 closed',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  UNIQUE KEY uk_part_buyer (part_id, buyer_id),
  INDEX idx_buyer_id (buyer_id),
  INDEX idx_seller_id (seller_id)
);
```

### 8.17 reports

```sql
CREATE TABLE reports (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  reporter_id BIGINT NOT NULL,
  target_type TINYINT NOT NULL COMMENT '1 post, 2 comment, 3 part, 4 user',
  target_id BIGINT NOT NULL,
  reason_type TINYINT NOT NULL,
  reason_text VARCHAR(500) NOT NULL DEFAULT '',
  status TINYINT NOT NULL DEFAULT 0 COMMENT '0 pending, 1 processed, 2 ignored',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  INDEX idx_target (target_type, target_id),
  INDEX idx_status_created_at (status, created_at)
);
```

### 8.18 audit_logs

```sql
CREATE TABLE audit_logs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  admin_id BIGINT NOT NULL,
  target_type TINYINT NOT NULL COMMENT '1 post, 2 comment, 3 part, 4 user, 5 report',
  target_id BIGINT NOT NULL,
  action VARCHAR(64) NOT NULL,
  remark VARCHAR(500) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL,
  INDEX idx_target (target_type, target_id)
);
```

### 8.19 admin_users

```sql
CREATE TABLE admin_users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(32) NOT NULL DEFAULT 'operator',
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

## 9. API 清单

### 9.1 认证

```text
POST /api/v1/auth/wechat-login
GET  /api/v1/auth/me
POST /api/v1/auth/logout
POST /api/v1/auth/bind-phone
```

### 9.2 用户和车库

```text
GET  /api/v1/users/me
PUT  /api/v1/users/me
GET  /api/v1/users/me/garages
POST /api/v1/users/me/garages
PUT  /api/v1/users/me/garages/{garageId}
DELETE /api/v1/users/me/garages/{garageId}
```

### 9.3 车型库

```text
GET /api/v1/vehicles/brands
GET /api/v1/vehicles/brands/{brandId}/series
GET /api/v1/vehicles/series/{seriesId}/models
GET /api/v1/vehicles/search
```

### 9.4 社区帖子

```text
GET  /api/v1/posts
GET  /api/v1/posts/{postId}
POST /api/v1/posts
PUT  /api/v1/posts/{postId}
DELETE /api/v1/posts/{postId}
POST /api/v1/posts/{postId}/like
DELETE /api/v1/posts/{postId}/like
POST /api/v1/posts/{postId}/favorite
DELETE /api/v1/posts/{postId}/favorite
GET  /api/v1/posts/{postId}/comments
POST /api/v1/posts/{postId}/comments
DELETE /api/v1/comments/{commentId}
```

### 9.5 二手件

```text
GET  /api/v1/parts
GET  /api/v1/parts/{partId}
POST /api/v1/parts
PUT  /api/v1/parts/{partId}
DELETE /api/v1/parts/{partId}
POST /api/v1/parts/{partId}/favorite
DELETE /api/v1/parts/{partId}/favorite
GET  /api/v1/part-categories
```

查询参数示例：

```text
GET /api/v1/parts?type=1&vehicleModelId=1001&categoryId=3&cityCode=310000&minPrice=1000&maxPrice=5000
```

### 9.6 交易意向

```text
POST /api/v1/parts/{partId}/intents
GET  /api/v1/intents
GET  /api/v1/intents/{intentId}
POST /api/v1/intents/{intentId}/close
```

### 9.7 举报

```text
POST /api/v1/reports
GET  /api/v1/reports/my
```

### 9.8 上传

```text
POST /api/v1/uploads/signature
POST /api/v1/uploads/complete
```

### 9.9 后台接口

```text
POST /api/v1/admin/auth/login
GET  /api/v1/admin/users
PUT  /api/v1/admin/users/{userId}/status
GET  /api/v1/admin/posts/review
GET  /api/v1/admin/posts/pending-review
POST /api/v1/admin/posts/{postId}/approve
POST /api/v1/admin/posts/{postId}/hide
POST /api/v1/admin/posts/{postId}/restore
GET  /api/v1/admin/parts/review
GET  /api/v1/admin/parts/pending-review
POST /api/v1/admin/parts/{partId}/approve
POST /api/v1/admin/parts/{partId}/hide
POST /api/v1/admin/parts/{partId}/restore
GET  /api/v1/admin/reports
POST /api/v1/admin/reports/{reportId}/process
GET  /api/v1/admin/vehicles
POST /api/v1/admin/vehicles/brands
POST /api/v1/admin/vehicles/series
POST /api/v1/admin/vehicles/models
GET  /api/v1/admin/categories
POST /api/v1/admin/categories
PUT  /api/v1/admin/categories/{categoryId}
```

## 10. 关键接口示例

### 10.1 发布二手件

```json
{
  "type": 1,
  "categoryId": 3,
  "title": "飞度 GK5 绞牙避震",
  "brand": "BC Racing",
  "model": "BR Series",
  "conditionLevel": 3,
  "price": 2800,
  "cityCode": "310000",
  "cityName": "上海",
  "description": "使用一年，正常拆车，无漏油，支持同城看货。",
  "fitments": [
    {
      "vehicleModelId": 1001,
      "note": "GK5 适配"
    }
  ],
  "images": [
    "https://cdn.example.com/parts/1.jpg"
  ]
}
```

返回示例：

```json
{
  "id": 12345,
  "status": 1
}
```

### 10.2 发起交易意向

```json
{
  "message": "我想要这套避震，准备按帖子里的联系方式沟通。"
}
```

返回示例：

```json
{
  "intentId": 9001,
  "status": 1
}
```

## 11. 审核和风控

### 11.1 内容审核

第一版采用混合审核策略。普通内容默认直接展示，系统和后台再进行巡检；命中高风险关键词的内容进入待审核，审核通过后才展示：

- 帖子
- 评论
- 二手件

MVP 建议采用：

- 敏感词规则
- 图片基础审核
- 用户举报
- 人工后台审核
- 高风险关键词命中后先审后发

### 11.2 高风险内容

重点拦截：

- 假冒品牌件
- 事故件隐瞒
- 安全隐患件
- 明确鼓励违法上路的改装内容
- 拆三元、炸街、遮挡号牌等违规导向内容
- 诈骗联系方式
- 引导私下高风险交易

说明：帖子和评论允许出现手机号、微信号，但仍需要识别诈骗、广告刷屏、诱导站外高风险交易和恶意导流。

### 11.3 平台提示

二手件详情页建议显示：

```text
平台仅提供信息撮合，不参与支付、担保、物流和售后。请优先同城验货，确认适配和成色后再交易。
```

## 12. 迭代计划

### 12.1 第 1 阶段：2-4 周

- 小程序基础框架
- Go API 基础框架
- 微信登录
- 手机号绑定
- 用户车库
- 帖子发布和列表
- 二手件发布和列表
- 图片上传
- 基础后台后审

### 12.2 第 2 阶段：4-6 周

- 评论、点赞、收藏
- 交易意向
- 举报
- 车型筛选
- 类目筛选
- 我的发布和我的收藏

### 12.3 第 3 阶段：6-8 周

- 车型库完善
- 搜索优化
- 内容推荐规则
- 用户等级
- 认证玩家
- 改装案例沉淀

## 13. 已确认决策

- 重点车型：第一版不限制，也不做官方推荐。
- 首发城市：第一版不限制。
- 联系方式：手机号与微信号可以在帖子和评论中出现。
- 改装店：第一版隐藏，只做普通用户。
- 内容审核：普通内容先发后审，高风险关键词命中内容审核后展示。
- 手机号绑定：第一版需要手机号绑定。
- 站内信/站内留言：第一版不做。
- 产品名称：中文名「撸车记」，英文名 retrofit。

## 14. 后续需要确认

- 高风险关键词的初始词库。
- 手机号绑定放在首次登录后立即要求，还是发布内容前要求。
- 是否允许用户在二手件详情页单独填写联系微信号。
