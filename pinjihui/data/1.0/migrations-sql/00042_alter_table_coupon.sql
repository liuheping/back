-- +goose Up
ALTER TABLE comments ALTER COLUMN is_show SET DEFAULT TRUE ;
-- +goose Down
ALTER TABLE comments ALTER COLUMN is_show SET DEFAULT FALSE ;
