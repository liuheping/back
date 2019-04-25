-- +goose Up
CREATE TABLE ads (
  id varchar(32) NOT NULL constraint ads_pkey primary key,
  image varchar(255) NOT NULL,
  link varchar(255),
  merchant_id varchar(32),
  position varchar(64) NOT NULL,
  sort int NOT NULL default 255,
  is_show bool not null default true
);
-- +goose Down
DROP TABLE ads;
