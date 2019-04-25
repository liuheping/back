-- +goose Up
ALTER TYPE shipping_address
  DROP ATTRIBUTE address;
ALTER TYPE shipping_address
  add attribute area_id integer;
ALTER TYPE shipping_address
  add attribute region_name varchar(255);
ALTER TYPE shipping_address
  add attribute address varchar(255);
-- +goose Down
ALTER TYPE shipping_address
  ADD ATTRIBUTE address address;
ALTER TYPE shipping_address
  DROP attribute area_id ;
ALTER TYPE shipping_address
  DROP attribute region_name ;
ALTER TYPE shipping_address
  DROP attribute address ;
