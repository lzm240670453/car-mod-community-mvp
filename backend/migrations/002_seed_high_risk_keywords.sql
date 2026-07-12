INSERT INTO high_risk_keywords (keyword, action, enabled) VALUES
  ('拆三元', 1, 1),
  ('炸街', 1, 1),
  ('遮挡号牌', 1, 1),
  ('套牌', 1, 1),
  ('假冒', 1, 1),
  ('事故件', 1, 1),
  ('无损刷程序', 1, 1)
ON DUPLICATE KEY UPDATE
  action = VALUES(action),
  enabled = VALUES(enabled);
