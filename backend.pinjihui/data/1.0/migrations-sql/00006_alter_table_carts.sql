-- +goose Up
alter table carts add column product_spec_id varchar(32);
ALTER TABLE public.carts
  ADD CONSTRAINT carts_product_specs_id_fk
FOREIGN KEY (product_spec_id) REFERENCES product_specs (id);

-- +goose Down
alter table carts drop column product_spec_id;
