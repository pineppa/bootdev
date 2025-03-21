-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at,
email, hashed_password, token, refresh_token, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    FALSE
)
RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE LOWER(name) = LOWER($1);

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByToken :one
SELECT * FROM users WHERE token = $1;

-- name: UpdateRedStatus :exec
UPDATE users
SET is_chirpy_red = $1 WHERE id = $2;

-- name: UpdateEmailPassword :exec
UPDATE users
SET email = $1, hashed_password = $2
WHERE token = $3;

-- name: SetTokenByID :exec
UPDATE users
SET token = $2 WHERE id = $1;

-- name: SetRefreshTokenByID :exec
UPDATE users
SET refresh_token = $2 WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;