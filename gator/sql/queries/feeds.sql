-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT p.* 
FROM posts p
JOIN feed_follows ff ON p.feed_id = ff.feed_id
JOIN users u ON ff.user_id = u.id
WHERE u.name = $1
ORDER BY p.published_at DESC
LIMIT $2;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1;

-- name: GetFeedFromId :one
SELECT * FROM feeds
WHERE id = $1;

-- name: GetAllFeeds :many
SELECT f.name, f.url, u.name AS username FROM feeds f
LEFT JOIN users u
ON u.id = f.user_id;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
    last_fetched_at = now(),
    updated_at = now()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;

-- name: DeleteAllFeeds :exec
DELETE FROM feeds;
DELETE FROM posts;
