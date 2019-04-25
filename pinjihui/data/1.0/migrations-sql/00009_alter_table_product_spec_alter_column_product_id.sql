-- +goose Up
ALTER TABLE product_specs ALTER COLUMN product_id SET NOT NULL;
-- +goose Down
alter table product_specs alter column product_id drop not null ;
