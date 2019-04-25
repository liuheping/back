-- +goose Up
ALTER TABLE products ADD COLUMN shipping_fee numeric;
-- +goose Down
ALTER TABLE products DROP COLUMN shipping_fee;
