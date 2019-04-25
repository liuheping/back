-- +goose Up
ALTER TABLE merchant_profiles ADD column waiters VARCHAR[];
-- +goose Down
ALTER TABLE merchant_profiles DROP column waiters;
