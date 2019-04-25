-- +goose Up
ALTER TABLE comments ADD COLUMN updated_at timestamp NOT NULL default current_timestamp;
CREATE TRIGGER update_time BEFORE UPDATE ON comments FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();
CREATE TRIGGER update_time BEFORE UPDATE ON cash_requests FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();
CREATE TRIGGER update_time BEFORE UPDATE ON coupons FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();
-- +goose Down
ALTER TABLE comments DROP COLUMN updated_at;
DROP TRIGGER update_time ON comments;
DROP TRIGGER update_time ON cash_requests;
DROP TRIGGER update_time ON coupons;
