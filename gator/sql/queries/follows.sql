-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    iff.*, 
    u.name AS user_name, 
    f.name AS feed_name
FROM inserted_feed_follow iff
LEFT JOIN users u ON u.id = iff.user_id
LEFT JOIN feeds f ON f.id = iff.feed_id;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows;

-- name: GetFeedFollowsForUser :many
SELECT u.name AS username, f.name AS feedname FROM feed_follows ff
JOIN users u ON u.id = ff.user_id
JOIN feeds f ON f.id = ff.feed_id
WHERE @username::TEXT = u.name;

-- name: GetFeedFollowsForUserFeedPair :one
SELECT * FROM feed_follows ff
JOIN users u ON u.id = ff.user_id
JOIN feeds f ON f.id = ff.feed_id
WHERE @username::TEXT = u.name AND @feedurl::TEXT = f.url;

-- name: UnfollowFeedFollow :exec
DELETE FROM feed_follows ff
USING users u, feeds f
WHERE ff.user_id = u.id
AND ff.feed_id = f.id
AND u.name = @username::TEXT
AND f.url = @feedurl::TEXT;

-- name: DeleteAllFeedFollows :exec
DELETE FROM feed_follows;
