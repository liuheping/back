-- +goose Up
-- 这个语句只能手动执行
ALTER TYPE coupon_type ADD VALUE 'for_inviter' AFTER 'for_newer';
-- +goose Down
--
