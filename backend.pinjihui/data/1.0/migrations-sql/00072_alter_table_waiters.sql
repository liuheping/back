-- +goose Up
ALTER TABLE public.waiters ALTER COLUMN merchant_id DROP NOT NULL;

-- +goose Down
ALTER TABLE public.waiters ALTER COLUMN merchant_id SET NOT null;
