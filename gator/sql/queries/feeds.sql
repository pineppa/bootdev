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

-- name: DeleteAllFeeds :exec
DELETE FROM feeds;
