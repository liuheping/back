-- +goose Up
ALTER TABLE comments ADD COLUMN is_anonymous bool NOT NULL DEFAULT FALSE;
ALTER TABLE comments ALTER COLUMN is_show SET DEFAULT TRUE ;
CREATE UNIQUE INDEX order_products_order_id_product_id_uindex ON public.order_products (order_id, product_id);
ALTER TABLE public.comments
  ADD CONSTRAINT comments_order_products_order_id_product_id_fk
FOREIGN KEY (order_id, product_id) REFERENCES order_products (order_id, product_id);
-- +goose Down
ALTER TABLE comments DROP COLUMN is_anonymous;
ALTER TABLE comments ALTER COLUMN is_show SET DEFAULT FALSE ;
DROP INDEX order_products_order_id_product_id_uindex;
ALTER TABLE comments DROP CONSTRAINT comments_order_products_order_id_product_id_fk;
