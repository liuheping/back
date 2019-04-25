-- +goose Up
ALTER TABLE orders ADD offer_amount numeric NOT NULL DEFAULT 0.0;
ALTER TABLE orders ALTER COLUMN used_coupon TYPE varchar[] USING used_coupon::varchar[];
-- +goose Down
ALTER TABLE orders DROP offer_amount;
ALTER TABLE orders ALTER COLUMN used_coupon TYPE varchar USING used_coupon::varchar;
