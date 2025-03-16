-- +goose Up
CREATE TABLE users (
    id          UUID PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    name        TEXT UNIQUE
);


-- +goose Down
DROP TABLE users;