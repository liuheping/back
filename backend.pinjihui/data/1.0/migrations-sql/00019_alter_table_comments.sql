-- +goose Up
ALTER TABLE comments ADD COLUMN shipping_rank SMALLINT;
ALTER TABLE comments ADD COLUMN service_rank SMALLINT;
ALTER TABLE comments ADD COLUMN images varchar[];
-- +goose Down
ALTER TABLE comments DROP COLUMN shipping_rank ;
ALTER TABLE comments DROP COLUMN service_rank ;
ALTER TABLE comments DROP COLUMN images;