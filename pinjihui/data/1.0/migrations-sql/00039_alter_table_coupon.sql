-- +goose Up
ALTER TABLE coupons ALTER COLUMN start_at TYPE date;
ALTER TABLE coupons ALTER COLUMN expired_at TYPE date;
ALTER TABLE user_coupons ALTER COLUMN start_at TYPE date;
ALTER TABLE user_coupons ALTER COLUMN expired_at TYPE date;
ALTER TABLE coupons ADD COLUMN merchant_id varchar(32) NULL;
ALTER TABLE user_coupons ADD COLUMN merchant_id varchar(32) NULL;

ALTER TABLE public.user_coupons
  ADD CONSTRAINT user_coupons_merchant_profiles_user_id_fk
FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);

ALTER TABLE coupons
  ADD CONSTRAINT coupons_merchant_profiles_user_id_fk
FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);

-- +goose Down
ALTER TABLE coupons ALTER COLUMN start_at TYPE timestamp;
ALTER TABLE user_coupons ALTER COLUMN start_at TYPE timestamp;
ALTER TABLE coupons ALTER COLUMN expired_at TYPE timestamp;
ALTER TABLE user_coupons ALTER COLUMN expired_at TYPE timestamp;
ALTER TABLE coupons DROP COLUMN merchant_id;
ALTER TABLE user_coupons DROP COLUMN merchant_id;
ALTER TABLE coupons DROP CONSTRAINT coupons_merchant_profiles_user_id_fk;
ALTER TABLE user_coupons DROP CONSTRAINT user_coupons_merchant_profiles_user_id_fk;
