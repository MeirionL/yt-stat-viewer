-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, channel_id, channel_name, access_token, refresh_token)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByChannelID :one
SELECT * FROM users WHERE channel_id = $1;

-- name: GetUserByChannelName :one
SELECT * FROM users WHERE channel_name = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET  updated_at = $2, channel_name = $3, access_token = $4, refresh_token = $5
WHERE id = $1
RETURNING *;

-- name: UpdateTokens :exec
UPDATE users
SET  updated_at = $2, access_token = $3, refresh_token = $4
WHERE id = $1
RETURNING *;

-- name: UpdateAccessToken :one
UPDATE users
SET  updated_at = $2, access_token = $3
WHERE id = $1
RETURNING *;

-- name: UpdateRefreshToken :one
UPDATE users
SET  updated_at = $2, refresh_token = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1 RETURNING channel_name;