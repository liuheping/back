-- +goose Up
CREATE EXTENSION cube;
CREATE EXTENSION earthdistance;
alter table merchant_profiles add column lat double precision null;
alter table merchant_profiles add column lng double precision null;
-- +goose Down
