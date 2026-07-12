CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  openid VARCHAR(64) NOT NULL UNIQUE,
  unionid VARCHAR(64) NULL,
  nickname VARCHAR(64) NOT NULL DEFAULT '',
  avatar_url VARCHAR(512) NOT NULL DEFAULT '',
  phone VARCHAR(32) NOT NULL DEFAULT '',
  phone_bound_at DATETIME NULL,
  status TINYINT NOT NULL DEFAULT 1 COMMENT '1 normal, 2 banned',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS vehicle_brands (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  initial CHAR(1) NOT NULL DEFAULT '',
  logo_url VARCHAR(512) NOT NULL DEFAULT '',
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_name (name),
  INDEX idx_initial_sort (initial, sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS vehicle_series (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  brand_id BIGINT NOT NULL,
  name VARCHAR(64) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_brand_id (brand_id),
  UNIQUE KEY uk_brand_series (brand_id, name),
  CONSTRAINT fk_vehicle_series_brand FOREIGN KEY (brand_id) REFERENCES vehicle_brands(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS vehicle_models (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  series_id BIGINT NOT NULL,
  name VARCHAR(128) NOT NULL,
  year SMALLINT NULL,
  generation VARCHAR(64) NOT NULL DEFAULT '',
  engine VARCHAR(64) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_series_id (series_id),
  INDEX idx_name (name),
  CONSTRAINT fk_vehicle_models_series FOREIGN KEY (series_id) REFERENCES vehicle_series(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS user_garages (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  vehicle_model_id BIGINT NOT NULL,
  year SMALLINT NULL,
  nickname VARCHAR(64) NOT NULL DEFAULT '',
  description VARCHAR(512) NOT NULL DEFAULT '',
  is_primary TINYINT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_user_id (user_id),
  INDEX idx_vehicle_model_id (vehicle_model_id),
  CONSTRAINT fk_user_garages_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_user_garages_vehicle_model FOREIGN KEY (vehicle_model_id) REFERENCES vehicle_models(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS posts (
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
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_user_id (user_id),
  INDEX idx_vehicle_model_id (vehicle_model_id),
  INDEX idx_status_created_at (status, created_at),
  FULLTEXT KEY ft_posts_title_content (title, content),
  CONSTRAINT fk_posts_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_posts_garage FOREIGN KEY (garage_id) REFERENCES user_garages(id),
  CONSTRAINT fk_posts_vehicle_model FOREIGN KEY (vehicle_model_id) REFERENCES vehicle_models(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS post_images (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL,
  image_url VARCHAR(512) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_post_id (post_id),
  CONSTRAINT fk_post_images_post FOREIGN KEY (post_id) REFERENCES posts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS comments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  parent_id BIGINT NULL,
  content VARCHAR(1000) NOT NULL,
  status TINYINT NOT NULL DEFAULT 1 COMMENT '0 pending_review, 1 visible, 2 hidden, 3 deleted',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_post_id (post_id),
  INDEX idx_user_id (user_id),
  INDEX idx_parent_id (parent_id),
  CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES posts(id),
  CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id) REFERENCES comments(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS post_likes (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_post_user (post_id, user_id),
  INDEX idx_user_id (user_id),
  CONSTRAINT fk_post_likes_post FOREIGN KEY (post_id) REFERENCES posts(id),
  CONSTRAINT fk_post_likes_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS post_favorites (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_post_user (post_id, user_id),
  INDEX idx_user_id (user_id),
  CONSTRAINT fk_post_favorites_post FOREIGN KEY (post_id) REFERENCES posts(id),
  CONSTRAINT fk_post_favorites_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS part_categories (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  parent_id BIGINT NULL,
  name VARCHAR(64) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_parent_id (parent_id),
  UNIQUE KEY uk_parent_name (parent_id, name),
  CONSTRAINT fk_part_categories_parent FOREIGN KEY (parent_id) REFERENCES part_categories(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS parts (
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
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_user_id (user_id),
  INDEX idx_category_id (category_id),
  INDEX idx_status_created_at (status, created_at),
  INDEX idx_city_code (city_code),
  INDEX idx_type_status_created_at (type, status, created_at),
  FULLTEXT KEY ft_parts_title_desc_brand_model (title, description, brand, model),
  CONSTRAINT fk_parts_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_parts_category FOREIGN KEY (category_id) REFERENCES part_categories(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS part_images (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  part_id BIGINT NOT NULL,
  image_url VARCHAR(512) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_part_id (part_id),
  CONSTRAINT fk_part_images_part FOREIGN KEY (part_id) REFERENCES parts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS part_fitments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  part_id BIGINT NOT NULL,
  vehicle_model_id BIGINT NOT NULL,
  note VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_part_vehicle (part_id, vehicle_model_id),
  INDEX idx_vehicle_model_id (vehicle_model_id),
  CONSTRAINT fk_part_fitments_part FOREIGN KEY (part_id) REFERENCES parts(id),
  CONSTRAINT fk_part_fitments_vehicle_model FOREIGN KEY (vehicle_model_id) REFERENCES vehicle_models(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS part_favorites (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  part_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_part_user (part_id, user_id),
  INDEX idx_user_id (user_id),
  CONSTRAINT fk_part_favorites_part FOREIGN KEY (part_id) REFERENCES parts(id),
  CONSTRAINT fk_part_favorites_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS trade_intents (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  part_id BIGINT NOT NULL,
  buyer_id BIGINT NOT NULL,
  seller_id BIGINT NOT NULL,
  message VARCHAR(500) NOT NULL DEFAULT '',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '1 active, 2 cancelled, 3 closed',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_part_buyer (part_id, buyer_id),
  INDEX idx_buyer_id (buyer_id),
  INDEX idx_seller_id (seller_id),
  CONSTRAINT fk_trade_intents_part FOREIGN KEY (part_id) REFERENCES parts(id),
  CONSTRAINT fk_trade_intents_buyer FOREIGN KEY (buyer_id) REFERENCES users(id),
  CONSTRAINT fk_trade_intents_seller FOREIGN KEY (seller_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS reports (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  reporter_id BIGINT NOT NULL,
  target_type TINYINT NOT NULL COMMENT '1 post, 2 comment, 3 part, 4 user',
  target_id BIGINT NOT NULL,
  reason_type TINYINT NOT NULL,
  reason_text VARCHAR(500) NOT NULL DEFAULT '',
  status TINYINT NOT NULL DEFAULT 0 COMMENT '0 pending, 1 processed, 2 ignored',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_reporter_id (reporter_id),
  INDEX idx_target (target_type, target_id),
  INDEX idx_status_created_at (status, created_at),
  CONSTRAINT fk_reports_reporter FOREIGN KEY (reporter_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS admin_users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(32) NOT NULL DEFAULT 'operator',
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  admin_id BIGINT NOT NULL,
  target_type TINYINT NOT NULL COMMENT '1 post, 2 comment, 3 part, 4 user, 5 report',
  target_id BIGINT NOT NULL,
  action VARCHAR(64) NOT NULL,
  remark VARCHAR(500) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_admin_id (admin_id),
  INDEX idx_target (target_type, target_id),
  CONSTRAINT fk_audit_logs_admin FOREIGN KEY (admin_id) REFERENCES admin_users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS high_risk_keywords (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  keyword VARCHAR(128) NOT NULL,
  action TINYINT NOT NULL DEFAULT 1 COMMENT '1 pending_review, 2 hide',
  enabled TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_keyword (keyword)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO part_categories (name, sort_order) VALUES
  ('轮毂轮胎', 10),
  ('避震悬挂', 20),
  ('刹车', 30),
  ('排气', 40),
  ('进气', 50),
  ('外观套件', 60),
  ('内饰', 70),
  ('灯光', 80),
  ('电子设备', 90),
  ('其他', 100)
ON DUPLICATE KEY UPDATE
  sort_order = VALUES(sort_order);
