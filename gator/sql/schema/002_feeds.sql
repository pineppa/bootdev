-- +goose Up
CREATE TABLE feeds (
    id              UUID PRIMARY KEY,
    created_at      TIMESTAMP NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP NOT NULL DEFAULT now(),
    last_fetched_at TIMESTAMP DEFAULT now(),
    name            TEXT UNIQUE NOT NULL,
    url             TEXT UNIQUE NOT NULL,
    user_id         UUID NOT NULL,
    FOREIGN KEY     (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE posts (
    id              UUID PRIMARY KEY,
    created_at      TIMESTAMP NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP NOT NULL DEFAULT now(),
    title           TEXT,
    url             TEXT UNIQUE,
    description     TEXT,
    published_at    TIMESTAMP,
    feed_id         UUID,
    FOREIGN KEY     (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE feeds;
DROP TABLE posts;