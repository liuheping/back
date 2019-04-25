-- +goose Up
ALTER TABLE public.shipping_info ADD COLUMN images varchar[];
ALTER TABLE public.shipping_info ALTER COLUMN company DROP NOT NULL ;
ALTER TABLE public.shipping_info ALTER COLUMN delivery_number DROP NOT NULL ;
-- +goose Down
ALTER TABLE public.shipping_info DROP COLUMN images;
ALTER TABLE public.shipping_info ALTER COLUMN company SET NOT NULL ;
ALTER TABLE public.shipping_info ALTER COLUMN delivery_number SET NOT NULL ;
