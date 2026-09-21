-- +goose Up

CREATE TYPE USER_ROLE AS ENUM ('user', 'admin');

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT UUIDV7(),

  email VARCHAR(320) NOT NULL,
  username VARCHAR(32) NOT NULL,
  role USER_ROLE NOT NULL DEFAULT 'user',

  password VARCHAR(255) NOT NULL,
  avatar VARCHAR(2048),

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX users_email_lower_idx ON users (LOWER(email));
CREATE UNIQUE INDEX users_username_lower_idx ON users (LOWER(username));

CREATE TABLE user_sessions (
  id UUID PRIMARY KEY DEFAULT UUIDV7(),

  user_id UUID NOT NULL,
  token_hash VARCHAR(128) NOT NULL,
  device_info VARCHAR(512) NOT NULL,
  device_hash VARCHAR(128) NOT NULL,

  ip_address INET NOT NULL,
  last_ip_address INET NOT NULL,

  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_user_sessions_user
  FOREIGN KEY (user_id)
  REFERENCES users (id)
  ON DELETE CASCADE
);

CREATE INDEX idx_user_sessions_user_id
ON user_sessions (user_id);

CREATE INDEX idx_user_sessions_expires_at
ON user_sessions (expires_at);

-- +goose Down

DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS USER_ROLE;
