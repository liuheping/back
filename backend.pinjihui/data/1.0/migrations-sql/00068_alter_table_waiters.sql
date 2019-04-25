-- +goose Up
ALTER TABLE waiters RENAME COLUMN waiters_id TO waiter_id;
-- +goose Down
ALTER TABLE waiters RENAME COLUMN waiter_id TO waiters_id;
