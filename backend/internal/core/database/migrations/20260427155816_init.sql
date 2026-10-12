-- +goose Up

CREATE TABLE users (
    id TEXT PRIMARY KEY,

    email TEXT NOT NULL COLLATE nocase,
    username TEXT NOT NULL COLLATE nocase,
    role TEXT NOT NULL,

    password TEXT NOT NULL,
    avatar TEXT,

    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE UNIQUE INDEX users_email_nocase_idx ON users (email);
CREATE UNIQUE INDEX users_username_nocase_idx ON users (username);

-- +goose StatementBegin
CREATE TRIGGER trg_users_updated_at AFTER UPDATE ON users
FOR EACH ROW
WHEN (new.updated_at = old.updated_at)
BEGIN
    UPDATE users
    SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
    WHERE id = old.id;
END;
-- +goose StatementEnd

CREATE TABLE user_sessions (
    id TEXT PRIMARY KEY,

    user_id TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    device_info TEXT NOT NULL,
    device_hash TEXT NOT NULL,

    ip_address TEXT NOT NULL,
    last_ip_address TEXT NOT NULL,

    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),

    CONSTRAINT fk_user_sessions_user
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions (user_id);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions (expires_at);

-- +goose Down

DROP TRIGGER IF EXISTS trg_users_updated_at;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;
