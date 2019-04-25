-- +goose Up
ALTER TABLE public.comments
  ADD CONSTRAINT comments_rel_merchants_products_product_id_merchant_id_fk
FOREIGN KEY (product_id, merchant_id) REFERENCES rel_merchants_products (product_id, merchant_id);
-- +goose Down
ALTER TABLE comments DROP constraint comments_rel_merchants_products_product_id_merchant_id_fk;
