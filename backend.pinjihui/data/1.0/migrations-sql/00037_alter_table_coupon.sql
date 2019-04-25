-- +goose Up
ALTER TABLE coupons ADD COLUMN start_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE user_coupons ADD COLUMN start_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- +goose Down
ALTER TABLE coupons DROP COLUMN start_at;
ALTER TABLE user_coupons DROP COLUMN start_at;
