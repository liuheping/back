-- +goose Up
ALTER TABLE shipping_methods ADD COLUMN enable_for_platform bool NOT NULL DEFAULT TRUE;
UPDATE shipping_methods SET enable_for_platform=FALSE WHERE name='上门取货(自提)';
-- +goose Down
ALTER TABLE shipping_methods DROP COLUMN enable_for_platform;
