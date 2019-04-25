-- +goose Up
ALTER TABLE public.order_products
  ADD CONSTRAINT order_products_orders_id_fk
FOREIGN KEY (order_id) REFERENCES orders (id);
ALTER TABLE public.order_products
  ADD CONSTRAINT order_products_products_id_fk
FOREIGN KEY (product_id) REFERENCES products (id);
-- +goose Down
ALTER TABLE order_products DROP CONSTRAINT order_products_orders_id_fk;
ALTER TABLE order_products DROP CONSTRAINT order_products_products_id_fk;

