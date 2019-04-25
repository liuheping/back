-- +goose Up
ALTER TABLE order_products ADD column shipping_fee numeric;
-- +goose Down
ALTER TABLE order_products DROP column shipping_fee;
