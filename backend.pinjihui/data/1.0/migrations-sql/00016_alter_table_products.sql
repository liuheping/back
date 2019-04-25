-- +goose Up
CREATE TYPE product_type AS ENUM ('simple', 'configure');
ALTER TABLE products ADD COLUMN type product_type not null default 'simple';
-- +goose Down
DROP TYPE product_type;
ALTER TABLE products DROP COLUMN type;
