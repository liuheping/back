-- +goose Up
ALTER TABLE shipping_methods ADD COLUMN enabled bool NOT NULL DEFAULT TRUE;
-- +goose Down
ALTER TABLE shipping_methods DROP COLUMN enabled;
