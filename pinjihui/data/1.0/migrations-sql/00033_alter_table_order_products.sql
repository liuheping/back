-- +goose Up
ALTER TABLE order_products
  ADD COLUMN spec_1_name varchar(255);
ALTER TABLE order_products
  ADD COLUMN spec_2_name varchar(255);
ALTER TABLE order_products
  ADD COLUMN spec_1 varchar(255);
ALTER TABLE order_products
  ADD COLUMN spec_2 varchar(255);
-- +goose Down
ALTER TABLE order_products
  DROP spec_1_name;
ALTER TABLE order_products
  DROP spec_2_name;
ALTER TABLE order_products
  DROP spec_1;
ALTER TABLE order_products
  DROP spec_2;

