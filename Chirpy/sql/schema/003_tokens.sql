-- +goose Up
CREATE TABLE refresh_tokens (
    token       TEXT PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expired_at  TIMESTAMP DEFAULT NULL,
    revoked_at  TIMESTAMP DEFAULT NULL,
    expires     TIMESTAMP
);

-- +goose Down
DROP TABLE refresh_tokens;