-- +goose Up
CREATE TABLE feed_follows (
    id          UUID PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    user_id     UUID NOT NULL,
    feed_id     UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT  user_feed UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;