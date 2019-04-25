-- +goose Up
ALTER TABLE public.products ADD merchant_id varchar(32) NULL;
ALTER TABLE public.products
  ADD CONSTRAINT products_merchant_profiles_user_id_fk
FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);
update products set merchant_id=(select merchant_id from rel_merchants_products where product_id=products.id limit 1);
alter table products alter column merchant_id set not null;
-- +goose Down
alter table products drop column merchant_id;
