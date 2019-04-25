-- +goose Up
ALTER TABLE order_products ADD COLUMN batch_price numeric NOT NULL DEFAULT 0;
ALTER TABLE order_products ADD COLUMN second_price numeric NOT NULL DEFAULT 0;
UPDATE order_products SET (batch_price, second_price)=(SELECT batch_price, second_price FROM products WHERE products.id=order_products.product_id);
ALTER TABLE order_products ALTER COLUMN batch_price DROP DEFAULT;
ALTER TABLE order_products ALTER COLUMN second_price DROP DEFAULT;
ALTER TABLE order_products ALTER COLUMN product_price TYPE numeric;
-- +goose Down
ALTER TABLE order_products DROP COLUMN batch_price ;
ALTER TABLE order_products DROP COLUMN second_price ;
