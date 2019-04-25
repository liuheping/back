-- +goose Up
ALTER TABLE public.waiters DROP CONSTRAINT waiters_merchant_profiles_user_id_fk;

-- +goose Down
ALTER TABLE public.waiters ADD CONSTRAINT waiters_merchant_profiles_user_id_fk FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);
