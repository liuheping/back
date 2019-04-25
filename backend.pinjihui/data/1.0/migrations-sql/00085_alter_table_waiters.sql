-- +goose Up
ALTER TABLE public.coupons ALTER COLUMN expired_at DROP NOT NULL ;
ALTER TABLE public.coupons ALTER COLUMN start_at DROP NOT NULL ;

-- +goose Down
ALTER TABLE public.coupons ALTER COLUMN expired_at SET NOT NULL ;
ALTER TABLE public.coupons ALTER COLUMN start_at SET NOT NULL ;