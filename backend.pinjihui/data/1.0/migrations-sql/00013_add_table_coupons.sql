-- +goose Up
create table coupons (
  id varchar(32) not null
  constraint coupons_pkey
    primary key ,
  user_id varchar(32) not null,
  description varchar(255) not null,
  value numeric not null,
  used bool NOT NULL DEFAULT FALSE,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp
);

CREATE TRIGGER update_time BEFORE UPDATE ON public.coupons FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();

comment on table coupons is '优惠券';
comment on column coupons.value is '面值';

-- +goose Down
drop table coupons;
