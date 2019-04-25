-- +goose Up
ALTER TABLE products
  ADD COLUMN spec_1 varchar(255);
ALTER TABLE products
  ADD COLUMN spec_2 varchar(255);
ALTER TABLE products
  ADD COLUMN parent_id varchar(32);
ALTER TABLE products
  DROP COLUMN retail_price;
ALTER TABLE rel_merchants_products
  DROP COLUMN spec_id;
-- +goose Down
ALTER TABLE products
  DROP COLUMN spec_1;
ALTER TABLE products
  DROP COLUMN spec_2;
ALTER TABLE products
  DROP COLUMN parent_id;
ALTER TABLE products
  ADD COLUMN retail_price numeric not null DEFAULT 0;
ALTER TABLE rel_merchants_products
  ADD COLUMN spec_id varchar(32)
  constraint rel_merchants_products_product_spec_id_fk
  references product_specs;
