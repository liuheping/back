-- +goose Up
ALTER TABLE public.product_images ADD COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- +goose Down
ALTER TABLE public.product_images DROP COLUMN created_at;



