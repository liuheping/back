-- +goose Up
create table spikes (
  id          varchar(32) not null
    constraint spikes_pkey
    primary key,
  product_id  varchar(32) NOT NULL,
  price       numeric     NOT NULL,
  start_at    timestamp   NOT NULL,
  expired_at  timestamp   NOT NULL,
  total_count int4        NOT NULL,
  merchant_id varchar(32) NOT NULL,
  buy_limit   int4        NOT NULL
);

CREATE INDEX spikes_product_id_index
  ON spikes (product_id);
ALTER TABLE spikes
  ADD CONSTRAINT spikes_products_id_fk
FOREIGN KEY (product_id) REFERENCES products (id);
ALTER TABLE spikes
  ADD CONSTRAINT spikes_merchant_profiles_user_id_fk
FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);

comment on TABLE spikes
IS '秒杀表';
comment on column spikes.total_count is '秒杀商品总数';
comment on column spikes.buy_limit is '用户限购数量';

CREATE TABLE order_spikes (
  order_product_id varchar(32) NOT NULL
    constraint order_spikes_pkey primary key ,
  spike_id varchar(32) NOT NULL
);
ALTER TABLE order_spikes
  ADD CONSTRAINT order_spikes_order_products_id_fk
FOREIGN KEY (order_product_id) REFERENCES order_products (id);
ALTER TABLE order_spikes
  ADD CONSTRAINT order_spikes_spikes_id_fk
FOREIGN KEY (spike_id) REFERENCES spikes (id);

comment on TABLE order_spikes
IS '订单秒杀记录';
-- +goose Down
drop table spikes;
drop TABLE order_spikes;
