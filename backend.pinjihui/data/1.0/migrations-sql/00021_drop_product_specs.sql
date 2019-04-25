-- +goose Up
ALTER TABLE carts DROP COLUMN product_spec_id;
DROP TABLE product_specs;
-- +goose Down
