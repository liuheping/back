-- +goose Up
UPDATE public.brands SET second_price_ratio=0.1,retail_price_ratio=0.3;
ALTER TABLE public.brands ALTER COLUMN second_price_ratio SET NOT NULL ;
ALTER TABLE public.brands ALTER COLUMN retail_price_ratio SET NOT NULL ;

-- +goose Down
ALTER TABLE public.brands ALTER COLUMN second_price_ratio DROP NOT NULL ;
ALTER TABLE public.brands ALTER COLUMN retail_price_ratio DROP NOT NULL ;