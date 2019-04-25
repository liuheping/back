-- +goose Up
ALTER TABLE merchant_profiles ADD column balance numeric;
ALTER TABLE cash_request RENAME TO cash_requests;
ALTER TABLE cash_requests ADD COLUMN merchant_id varchar(32);
ALTER TABLE cash_requests ADD COLUMN updated_at timestamp NOT NULL default current_timestamp;
ALTER TABLE cash_requests ADD COLUMN created_at timestamp NOT NULL default current_timestamp;
-- +goose Down
ALTER TABLE merchant_profiles DROP COLUMN balance;
ALTER TABLE cash_requests DROP COLUMN merchant_id;
ALTER TABLE cash_requests DROP column created_at;
ALTER TABLE cash_requests DROP COLUMN updated_at;
ALTER TABLE cash_requests RENAME TO cash_request;
