-- +goose Up
CREATE SEQUENCE users_invited_code_seq
  start with 100000;
ALTER TABLE users DROP COLUMN invite_code;
ALTER TABLE users
  ADD COLUMN invite_code integer NOT NULL DEFAULT nextval('users_invited_code_seq');
-- +goose Down
DROP SEQUENCE users_invited_code_seq;
ALTER TABLE users ALTER COLUMN invite_code TYPE varchar(32);
