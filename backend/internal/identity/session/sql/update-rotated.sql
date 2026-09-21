UPDATE user_sessions
SET token_hash = $1, last_ip_address = $2, expires_at = $3
WHERE id = $4;
