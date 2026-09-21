SELECT
  id,
  email,
  username,
  role,
  password,
  avatar,
  created_at,
  updated_at
FROM users
WHERE username = $1;
