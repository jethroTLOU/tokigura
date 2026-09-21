SELECT
  id,
  user_id,
  token_hash,
  device_info,
  device_hash,
  ip_address,
  last_ip_address,
  expires_at,
  created_at
FROM user_sessions
WHERE token_hash = $1;
