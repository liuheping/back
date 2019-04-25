-- +goose Up
ALTER TABLE wecharts ALTER COLUMN user_id SET NOT NULL;
-- +goose Down
ALTER TABLE wecharts ALTER COLUMN user_id DROP NOT NULL;
