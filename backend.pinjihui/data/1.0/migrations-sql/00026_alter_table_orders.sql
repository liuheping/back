-- +goose Up
ALTER TABLE orders ADD COLUMN used_coupon varchar(32) NULL;
ALTER TABLE orders ALTER COLUMN shipping_id DROP NOT NULL ;
ALTER TABLE orders ALTER COLUMN shipping_name DROP NOT NULL ;
-- +goose Down
ALTER TABLE orders DROP COLUMN used_coupon;
ALTER TABLE orders ALTER COLUMN shipping_id SET NOT NULL ;
ALTER TABLE orders ALTER COLUMN shipping_name SET NOT NULL ;
