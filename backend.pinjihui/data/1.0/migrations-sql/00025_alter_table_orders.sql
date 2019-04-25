-- +goose Up
ALTER TABLE orders RENAME COLUMN pay_id TO pay_code;
-- +goose Down
ALTER TABLE orders RENAME COLUMN pay_code TO pay_id;
