-- +goose Up
ALTER TABLE public.merchant_profiles DROP COLUMN waiters;
-- +goose Down
ALTER TABLE public.merchant_profiles ADD COLUMN waiters VARCHAR[];


