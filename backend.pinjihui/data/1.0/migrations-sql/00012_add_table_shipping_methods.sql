-- +goose Up
create table shipping_methods (
  id varchar(32) not null
  constraint shipping_methods_pkey
    primary key ,
  name varchar(32) not null,
  fee numeric not null
);

comment on table shipping_methods is '配送方式表';
comment on column shipping_methods.name is '配送方式名称';
comment on column shipping_methods.fee is '配送费用';

-- +goose Down
drop table shipping_methods;
