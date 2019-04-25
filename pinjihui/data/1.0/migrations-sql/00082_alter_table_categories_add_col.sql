-- +goose Up
ALTER TABLE categories ADD COLUMN is_common BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose Down
ALTER TABLE categories DROP COLUMN is_common;