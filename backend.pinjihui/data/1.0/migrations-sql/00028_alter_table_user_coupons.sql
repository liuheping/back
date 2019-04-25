-- +goose Up
ALTER TABLE user_coupons
  ADD COLUMN limit_amount numeric NULL;
ALTER TABLE coupons RENAME COLUMN limitamount TO limit_amount;
-- +goose Down
ALTER TABLE user_coupons
  DROP COLUMN limitAmount;
