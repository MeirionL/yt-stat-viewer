-- +goose Up
ALTER TABLE users ADD COLUMN access_token TEXT NOT NULL, refresh_token TEXT NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN access_token, refresh_token;