-- +goose Up
ALTER TABLE merchant_profiles ADD column logo VARCHAR(255);
ALTER TABLE merchant_profiles ADD column telephone VARCHAR(16);
-- +goose Down
ALTER TABLE merchant_profiles DROP column logo;
ALTER TABLE merchant_profiles DROP column telephone;
