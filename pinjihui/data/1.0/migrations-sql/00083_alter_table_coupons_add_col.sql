-- +goose Up
ALTER TABLE coupons ADD COLUMN validity_days INTEGER NULL ;
-- +goose Down
ALTER TABLE coupons DROP COLUMN validity_days;