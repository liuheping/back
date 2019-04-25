-- +goose Up
alter table products alter column retail_price set not null;

-- +goose Down
alter table products alter column retail_price drop not null;
