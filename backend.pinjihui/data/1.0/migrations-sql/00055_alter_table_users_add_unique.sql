-- +goose Up
CREATE UNIQUE INDEX users_invite_code_uindex ON public.users (invite_code);
-- +goose Down
DROP INDEX users_invite_code_uindex;
