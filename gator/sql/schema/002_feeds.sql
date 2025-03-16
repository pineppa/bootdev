-- +goose Up
CREATE TABLE feeds (
    id          UUID PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    name        TEXT UNIQUE NOT NULL,
    url         TEXT UNIQUE NOT NULL,
    user_id     UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE feeds;