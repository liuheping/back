-- +goose Up
ALTER TABLE public.rel_merchants_products ADD COLUMN is_sale BOOL NOT NULL DEFAULT TRUE;

-- +goose Down
ALTER TABLE public.rel_merchants_products DROP COLUMN is_sale;