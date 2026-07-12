CREATE TABLE IF NOT EXISTS site_messages (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  recipient_id BIGINT NOT NULL,
  actor_id BIGINT NULL,
  type TINYINT NOT NULL COMMENT '1 system, 2 trade, 3 interaction',
  title VARCHAR(120) NOT NULL,
  content VARCHAR(500) NOT NULL DEFAULT '',
  target_type TINYINT NOT NULL DEFAULT 0 COMMENT '0 none, 1 post, 2 comment, 3 part, 4 user',
  target_id BIGINT NOT NULL DEFAULT 0,
  read_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_recipient_created_at (recipient_id, created_at),
  INDEX idx_recipient_read_at (recipient_id, read_at),
  INDEX idx_target (target_type, target_id),
  CONSTRAINT fk_site_messages_recipient FOREIGN KEY (recipient_id) REFERENCES users(id),
  CONSTRAINT fk_site_messages_actor FOREIGN KEY (actor_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
