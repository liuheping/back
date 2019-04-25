-- +goose Up
ALTER TABLE cash_requests ALTER COLUMN merchant_id SET NOT NULL;
ALTER TABLE cash_requests
  ADD CONSTRAINT cash_requests_merchant_profiles_user_id_fk
FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);
-- +goose Down
ALTER TABLE cash_requests ALTER COLUMN merchant_id DROP NOT NULL;
ALTER TABLE cash_requests
  DROP CONSTRAINT cash_requests_merchant_profiles_user_id_fk;
