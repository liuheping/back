-- +goose Up
-- 这个语句只能手动执行
ALTER TYPE coupon_type ADD VALUE 'for_sharer';
ALTER TYPE coupon_type ADD VALUE 'for_be_sharer';
ALTER TYPE coupon_type ADD VALUE 'for_first_login';
-- +goose Down
--
