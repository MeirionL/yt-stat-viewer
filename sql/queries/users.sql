-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, platform, access_token, refresh_token)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

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

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;