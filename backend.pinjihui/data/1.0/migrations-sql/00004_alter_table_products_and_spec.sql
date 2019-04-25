-- +goose Up
alter table products add column spec_1_name varchar(255);
alter table products add column spec_2_name varchar(255);

CREATE UNIQUE INDEX product_specs_sku_uindex ON public.product_specs (sku);

-- +goose Down
alter table products drop column spec_1_name;
alter table products drop column spec_2_name;

drop index product_specs_sku_uindex;
