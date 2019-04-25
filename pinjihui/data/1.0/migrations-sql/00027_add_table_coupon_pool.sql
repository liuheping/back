-- +goose Up
ALTER TABLE coupons
  ADD COLUMN expired_at TIMESTAMP NOT NULL DEFAULT current_timestamp + INTERVAL '1 year';
ALTER TABLE coupons
  RENAME TO user_coupons;
ALTER INDEX coupons_pkey RENAME TO user_coupons_pkey;
create table coupons (
  id          varchar(32)                         not null
    constraint coupons_pkey
    primary key,
  description varchar(255)                        not null,
  value       numeric                             not null,
  created_at  timestamp default CURRENT_TIMESTAMP not null,
  updated_at  timestamp default CURRENT_TIMESTAMP not null,
  limitAmount numeric                             NULL,
  expired_at  timestamp                           NOT NULL
);
comment on table coupons
is '优惠券';
INSERT INTO coupons
VALUES ('coupons_for_newer', '新手优惠券', 50.00, current_timestamp, current_timestamp, 350.00, timestamp '2018-06-07 10:16:15' + INTERVAL '1 year');

-- +goose Down
drop table coupons;
ALTER INDEX coupons_pkey RENAME TO user_coupons_pkey;
ALTER TABLE user_coupons
  RENAME TO coupons;
ALTER TABLE coupons
  DROP COLUMN expired_at;
