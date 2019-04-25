-- +goose Up
ALTER TABLE orders ADD COLUMN inv_taxpayer_id varchar(255);
ALTER TABLE orders ADD COLUMN inv_url varchar(255);
ALTER TABLE orders DROP COLUMN how_oos;
ALTER TABLE orders DROP COLUMN inv_content;
ALTER TABLE orders ALTER COLUMN amount TYPE numeric;
ALTER TABLE orders ALTER COLUMN shipping_fee TYPE numeric;
ALTER TABLE orders ALTER COLUMN pay_fee TYPE numeric;
ALTER TABLE orders ALTER COLUMN money_paid TYPE numeric;
ALTER TABLE orders ALTER COLUMN order_amount TYPE numeric;
-- +goose Down
ALTER TABLE orders DROP COLUMN inv_taxpayer_id ;
ALTER TABLE orders DROP COLUMN inv_url;
ALTER TABLE orders ADD COLUMN how_oos how_oos;
ALTER TABLE orders ADD COLUMN inv_content varchar(255);
ALTER TABLE orders ALTER COLUMN amount TYPE money;
ALTER TABLE orders ALTER COLUMN shipping_fee TYPE money;
ALTER TABLE orders ALTER COLUMN pay_fee TYPE money;
ALTER TABLE orders ALTER COLUMN money_paid TYPE money;
ALTER TABLE orders ALTER COLUMN order_amount TYPE money;
