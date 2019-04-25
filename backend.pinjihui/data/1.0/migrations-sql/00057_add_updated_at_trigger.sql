-- +goose Up
CREATE TRIGGER update_time BEFORE UPDATE ON debit_card_info FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();
CREATE TRIGGER update_time BEFORE UPDATE ON products FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();
CREATE TRIGGER update_time BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();
-- +goose Down
DROP TRIGGER update_time ON debit_card_info;
DROP TRIGGER update_time ON products;
DROP TRIGGER update_time ON users;
