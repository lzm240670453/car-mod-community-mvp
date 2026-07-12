INSERT INTO admin_users (username, password_hash, role, status) VALUES
  ('admin', 'dev-password-not-checked', 'operator', 1)
ON DUPLICATE KEY UPDATE
  role = VALUES(role),
  status = VALUES(status);
