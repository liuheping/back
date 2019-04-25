-- +goose Up
CREATE TYPE coupon_type AS ENUM ('for_newer', 'simple');
ALTER TABLE coupons ADD COLUMN quantity int4 default 0 not null ;
ALTER TABLE coupons ADD COLUMN type coupon_type NOT NULL DEFAULT 'simple';
UPDATE coupons SET id='bbp4iguak7rhnps11ke0', type='for_newer' WHERE id='coupons_for_newer';
ALTER TABLE user_coupons ADD COLUMN type coupon_type NOT NULL DEFAULT 'simple';
-- +goose Down
ALTER TABLE coupons
  DROP COLUMN quantity;
ALTER TABLE coupons DROP COLUMN type;
UPDATE coupons SET id='coupons_for_newer' WHERE id='bbp4iguak7rhnps11ke0';
ALTER TABLE user_coupons DROP COLUMN type;