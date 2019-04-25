-- +goose Up
ALTER TABLE orders ADD COLUMN provider_income numeric ;
ALTER TABLE orders ADD COLUMN ally_income     numeric ;

UPDATE orders SET (provider_income, ally_income) = (SELECT SUM(batch_price), SUM(product_price-second_price) FROM order_products WHERE order_products.order_id=orders.id);
-- +goose Down
ALTER TABLE orders DROP COLUMN provider_income;
ALTER TABLE orders DROP COLUMN ally_income;
