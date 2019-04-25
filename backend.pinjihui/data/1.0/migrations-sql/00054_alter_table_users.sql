-- +goose Up
ALTER TABLE users ADD COLUMN invite_code VARCHAR(32);
ALTER TABLE users ADD COLUMN invited BOOLEAN NOT NULL DEFAULT FALSE ;
-- +goose Down
ALTER TABLE users DROP COLUMN invite_code ;
ALTER TABLE users DROP COLUMN invited ;
