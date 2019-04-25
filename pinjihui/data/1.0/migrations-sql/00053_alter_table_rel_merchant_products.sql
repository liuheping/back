-- +goose Up
ALTER TABLE rel_merchants_products ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ;
ALTER TABLE rel_merchants_products ADD COLUMN sales_volume INT NOT NULL DEFAULT 0;

UPDATE merchant_profiles SET company_image='{'||company_image||'}';
ALTER TABLE merchant_profiles ALTER COLUMN company_image TYPE VARCHAR[] USING company_image::VARCHAR[];
-- +goose Down
ALTER TABLE rel_merchants_products DROP COLUMN created_at ;
ALTER TABLE rel_merchants_products DROP COLUMN sales_volume ;

ALTER TABLE merchant_profiles ALTER COLUMN company_image TYPE VARCHAR USING company_image::VARCHAR;
