-- +goose Up
INSERT INTO public.users (id, name, mobile, password, type, email, created_at, updated_at, status, last_ip, last_login_time) VALUES ('platform', null, '13540069749', '$2a$10$fSzovH3pg1JizkoZCBKwje9CSqnUAInXa3vxEkruOozl0qUslmfMm', 'admin', null, '2018-05-28 11:11:05.990497', '2018-05-28 11:11:05.990497', 'normal', '::1', null) ON CONFLICT DO NOTHING ;
INSERT INTO public.merchant_profiles (user_id, social_id, rep_name, company_name, company_address, delivery_address, license_image, company_image, created_at, updated_at, lat, lng) VALUES ('platform', '510922199816746478', '杨林', '成都拼机惠', null, null, null, null, '2018-06-06 11:34:28.331799', '2018-06-06 12:45:04.803513', null, null) ON CONFLICT DO NOTHING ;
INSERT INTO rel_merchants_products (SELECT product_id, 'platform', stock, retail_price-1 FROM rel_merchants_products) ON CONFLICT DO NOTHING ;
-- +goose Down
DELETE FROM users WHERE id='platform';
DELETE FROM rel_merchants_products WHERE merchant_id='platform';