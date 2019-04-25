-- +goose Up
create table shipping_info (
  id              varchar(32) not null
    constraint shipping_info_pkey
    primary key,
  order_id varchar(32) NOT NULL ,
  company         varchar(32) NOT NULL,
  delivery_number varchar(32) NOT NULL
);

CREATE INDEX shipping_info_order_id_index ON shipping_info (order_id);
ALTER TABLE shipping_info
  ADD CONSTRAINT shipping_info_orders_id_fk
FOREIGN KEY (order_id) REFERENCES orders (id);

comment on TABLE shipping_info
IS '配送信息';
COMMENT ON COLUMN shipping_info.company IS '快递公司';
COMMENT ON COLUMN shipping_info.delivery_number IS '快递单号';

-- +goose Down
drop table shipping_info;
