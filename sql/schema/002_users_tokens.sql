-- +goose Up
ALTER TABLE users ADD COLUMN access_token TEXT NOT NULL;
ALTER TABLE users ADD COLUMN refresh_token TEXT NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN access_token;
ALTER TABLE users DROP COLUMN refresh_token;
