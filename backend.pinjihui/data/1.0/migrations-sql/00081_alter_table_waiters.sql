-- +goose Up
ALTER TABLE public.products ALTER COLUMN category_id SET NOT NULL ;
ALTER TABLE public.products ALTER COLUMN brand_id SET NOT NULL ;

-- +goose Down
ALTER TABLE public.products ALTER COLUMN category_id DROP NOT NULL ;
ALTER TABLE public.products ALTER COLUMN brand_id DROP NOT NULL ;