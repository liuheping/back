-- +goose Up
ALTER TABLE attribute_sets ALTER COLUMN merchant_id DROP NOT NULL ;
-- +goose Down
ALTER TABLE attribute_sets ALTER COLUMN merchant_id SET NOT NULL ;
