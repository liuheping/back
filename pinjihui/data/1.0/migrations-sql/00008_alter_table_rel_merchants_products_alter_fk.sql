-- +goose Up
ALTER TABLE public.rel_merchants_products DROP CONSTRAINT rel_merchants_products_users_id_fk;
ALTER TABLE public.rel_merchants_products
  ADD CONSTRAINT rel_merchants_products_merchant_profiles_user_id_fk
FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);
-- +goose Down
ALTER TABLE public.rel_merchants_products DROP CONSTRAINT rel_merchants_products_merchant_profiles_user_id_fk;
ALTER TABLE public.rel_merchants_products
  ADD CONSTRAINT rel_merchants_products_users_id_fk
FOREIGN KEY (merchant_id) REFERENCES users (id);
