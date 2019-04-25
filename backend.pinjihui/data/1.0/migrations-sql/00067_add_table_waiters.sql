-- +goose Up
CREATE TABLE waiters (
  id VARCHAR(32) NOT NULL CONSTRAINT waiters_pkey PRIMARY KEY,
  merchant_id VARCHAR(32) NOT NULL,
  mobile varchar(32) NOT NULL,
  name VARCHAR(32),
  waiters_id VARCHAR(32) ,
  checked BOOL NOT NULL DEFAULT FALSE
);
-- +goose Down
DROP TABLE waiters;
