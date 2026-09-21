INSERT INTO user_sessions (
  user_id, token_hash, device_info, device_hash,
  ip_address, last_ip_address, expires_at
) VALUES (
  $1, $2, $3, $4,
  $5, $6, $7
);
