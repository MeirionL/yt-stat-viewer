-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, api_key, oauth_token)
VALUES ($1, $2, $3, $4, encode(sha256(random()::text::bytea), 'hex'), $5)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: UpdateUser :one
UPDATE users
SET  updated_at = $2, name = $3, oauth_token = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;