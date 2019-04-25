-- +goose Up
ALTER TABLE public.categories
  ADD CONSTRAINT categories_categories_id_fk
FOREIGN KEY (parent_id) REFERENCES categories (id);

ALTER TABLE public.comments
  ADD CONSTRAINT comments_products_id_fk
FOREIGN KEY (product_id) REFERENCES products (id);
ALTER TABLE public.comments
  ADD CONSTRAINT comments_orders_id_fk
FOREIGN KEY (order_id) REFERENCES orders (id);
ALTER TABLE public.comments
  ADD CONSTRAINT comments_merchant_profiles_user_id_fk
FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);

ALTER TABLE public.merchant_profiles
  ADD CONSTRAINT merchant_profiles_users_id_fk
FOREIGN KEY (user_id) REFERENCES public.users (id);

ALTER TABLE public.orders
  ADD CONSTRAINT orders_users_id_fk
FOREIGN KEY (user_id) REFERENCES public.users (id);

CREATE UNIQUE INDEX payments_pay_code_uindex ON public.payments (pay_code);

ALTER TABLE public.orders
  ADD CONSTRAINT orders_payments_pay_code_fk
FOREIGN KEY (pay_code) REFERENCES payments (pay_code);

ALTER TABLE public.orders
  ADD CONSTRAINT orders_orders_id_fk
FOREIGN KEY (parent_id) REFERENCES orders (id);
ALTER TABLE public.orders
  ADD CONSTRAINT orders_merchant_profiles_user_id_fk
FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);

ALTER TABLE public.product_images
  ADD CONSTRAINT product_images_products_id_fk
FOREIGN KEY (product_id) REFERENCES products (id);

ALTER TABLE public.user_coupons
  ADD CONSTRAINT user_coupons_users_id_fk
FOREIGN KEY (user_id) REFERENCES public.users (id);

ALTER TABLE public.wecharts
  ADD CONSTRAINT wecharts_users_id_fk
FOREIGN KEY (user_id) REFERENCES public.users (id);

-- +goose Down
ALTER TABLE categories DROP CONSTRAINT categories_categories_id_fk;
ALTER TABLE comments DROP CONSTRAINT comments_products_id_fk;
ALTER TABLE comments DROP CONSTRAINT comments_orders_id_fk;
ALTER TABLE comments DROP CONSTRAINT comments_merchant_profiles_user_id_fk;
ALTER TABLE merchant_profiles DROP CONSTRAINT merchant_profiles_users_id_fk;
ALTER TABLE orders DROP CONSTRAINT orders_users_id_fk;
DROP INDEX payments_pay_code_uindex;
ALTER TABLE orders DROP CONSTRAINT orders_payments_pay_code_fk;
ALTER TABLE orders DROP CONSTRAINT orders_orders_id_fk;
ALTER TABLE orders DROP CONSTRAINT orders_merchant_profiles_user_id_fk;
ALTER TABLE product_images DROP CONSTRAINT product_images_products_id_fk;
ALTER TABLE user_coupons DROP CONSTRAINT user_coupons_users_id_fk;
ALTER TABLE wecharts DROP CONSTRAINT wecharts_users_id_fk;
