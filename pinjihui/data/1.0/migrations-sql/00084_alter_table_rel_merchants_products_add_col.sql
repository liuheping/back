-- +goose Up
ALTER TABLE rel_merchants_products ADD COLUMN view_volume INTEGER NOT NULL DEFAULT 0;
-- +goose Down
ALTER TABLE rel_merchants_products DROP COLUMN view_volume;