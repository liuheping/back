-- +goose Up
ALTER TABLE favorites RENAME COLUMN object_id to merchant_id;
ALTER TABLE favorites ADD COLUMN product_id varchar(32) null;
ALTER TABLE favorites DROP COLUMN type;

ALTER TABLE public.favorites
  ADD CONSTRAINT favorites_merchant_profiles_user_id_fk
FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);
ALTER TABLE public.favorites
  ADD CONSTRAINT favorites_products_id_fk
FOREIGN KEY (product_id) REFERENCES products (id);
ALTER TABLE public.favorites
  ADD CONSTRAINT favorites_users_id_fk
FOREIGN KEY (user_id) REFERENCES public.users (id);
CREATE UNIQUE INDEX favorites_user_id_merchant_id_product_id_uindex ON public.favorites (user_id, merchant_id, product_id);
-- +goose Down
ALTER TABLE favorites RENAME COLUMN merchant_id to object_id;
ALTER TABLE favorites DROP COLUMN product_id;
ALTER TABLE favorites ADD COLUMN type favorites NOT NULL DEFAULT 'product';

ALTER TABLE favorites DROP CONSTRAINT favorites_merchant_profiles_user_id_fk;
ALTER TABLE favorites DROP CONSTRAINT favorites_products_id_fk;
ALTER TABLE favorites DROP CONSTRAINT favorites_users_id_fk;
DROP INDEX favorites_user_id_merchant_id_product_id_uindex;