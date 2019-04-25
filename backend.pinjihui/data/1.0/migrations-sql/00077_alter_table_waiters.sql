-- +goose Up
ALTER TABLE public.waiters DROP COLUMN handled;

-- +goose Down
ALTER TABLE public.waiters ADD COLUMN handled Bool NOT NULL DEFAULT FALSE;



