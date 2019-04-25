-- +goose Up
ALTER TABLE rel_merchants_products ADD COLUMN origin_price numeric;
-- +goose Down
ALTER TABLE rel_merchants_products DROP COLUMN origin_price;