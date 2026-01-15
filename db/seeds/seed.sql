-- auth4x seed data
INSERT INTO users (id, username, password_hash)
VALUES
  ('00000000-0000-0000-0000-000000000001', 'admin', 'bcrypt:$2a$10$exampleadminhash'),
  ('00000000-0000-0000-0000-000000000002', 'demo', 'bcrypt:$2a$10$exampledemohash');

INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
VALUES
  ('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'sha256:exampletokenhash', CURRENT_TIMESTAMP + INTERVAL '30 day');
