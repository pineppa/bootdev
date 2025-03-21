-- name: CreateToken :one
INSERT INTO refresh_tokens (
    token, created_at, updated_at, 
    user_id, expired_at, revoked_at, expires)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NULL,
    NULL,
    NOW() + INTERVAL '60 day'
)
RETURNING *;

-- name: GetAllTokens :many
SELECT * FROM refresh_tokens;

-- name: GetTokenByToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: GetTokenByUserId :one
SELECT * FROM refresh_tokens WHERE user_id = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET expired_at = NOW(), expires = NOW(), revoked_at = NOW()
WHERE token = $1;

-- name: DeleteAllTokens :exec
DELETE FROM refresh_tokens;