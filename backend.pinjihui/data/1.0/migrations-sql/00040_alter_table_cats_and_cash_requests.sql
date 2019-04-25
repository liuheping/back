-- +goose Up
ALTER TABLE cash_requests ALTER COLUMN amount TYPE numeric;
ALTER TABLE categories DROP COLUMN deleted;
-- +goose Down
ALTER TABLE cash_requests ALTER COLUMN amount TYPE money;
ALTER TABLE categories ADD COLUMN deleted boolean DEFAULT FALSE NOT NULL ;
