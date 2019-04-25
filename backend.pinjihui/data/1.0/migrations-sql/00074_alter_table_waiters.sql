-- +goose Up
ALTER TABLE public.brands ADD COLUMN second_price_ratio NUMERIC(100,5);
ALTER TABLE public.brands ADD COLUMN retail_price_ratio NUMERIC(100,5);
-- +goose Down
ALTER TABLE public.brands DROP COLUMN second_price_ratio;
ALTER TABLE public.brands DROP COLUMN retail_price_ratio;


