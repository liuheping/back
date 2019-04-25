-- +goose Up
update rel_merchants_products set retail_price=(select retail_price from products where rel_merchants_products.product_id=products.id limit 1);
ALTER TABLE public.rel_merchants_products ALTER COLUMN retail_price SET NOT NULL;
-- +goose Down
ALTER TABLE public.rel_merchants_products ALTER COLUMN retail_price drop NOT NULL;
