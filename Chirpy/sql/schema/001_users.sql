-- +goose Up
CREATE TABLE users (
    id              UUID PRIMARY KEY,
    created_at      TIMESTAMP NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP NOT NULL DEFAULT now(),
    email           TEXT UNIQUE NOT NULL,
    hashed_password TEXT,
    expires_sec     INT,
    token           TEXT,
    refresh_token   TEXT,
    is_chirpy_red   BOOLEAN NOT NULL DEFAULT FALSE
);


-- +goose Down
DROP TABLE users;