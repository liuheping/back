-- +goose Up
create table product_specs (
  id varchar(32) not null
  constraint product_specs_pkey
    primary key ,
  product_id varchar(32)
  constraint product_specs_products_id_fk
    references products,
  spec_1 varchar(255) not null ,
  spec_2 varchar(255),

  batch_price      numeric                     not null,
  second_price     numeric                     not null,
  retail_price     numeric                     not null,
  stock int4 not null,
  sku varchar(32) not null
);

comment on table product_specs is '规格表';

alter table rel_merchants_products add column spec_id varchar(32)
  constraint rel_merchants_products_product_spec_id_fk
    references product_specs;
-- +goose Down
drop table product_specs;
alter table rel_merchants_products drop column spec_id;