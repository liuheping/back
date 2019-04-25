-- +goose Up
ALTER TABLE waiters ADD COLUMN deleted Bool NOT NULL DEFAULT FALSE ;
ALTER TABLE waiters ADD COLUMN handled Bool NOT NULL DEFAULT FALSE ;

-- +goose Down
ALTER TABLE waiters DROP COLUMN deleted;
ALTER TABLE waiters DROP COLUMN handled;