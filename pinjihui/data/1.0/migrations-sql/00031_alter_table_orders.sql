-- +goose Up
ALTER TABLE orders DROP COLUMN shipping_status;
ALTER TABLE orders DROP COLUMN pay_status;
ALTER TABLE orders DROP COLUMN order_status;
DROP type order_status;
CREATE type order_status AS ENUM ('unpaid', 'paid', 'shipped', 'finish', 'cancelled', 'invalid', 'returned');
ALTER TABLE orders ADD COLUMN status order_status NOT NULL DEFAULT 'unpaid';
-- +goose Down
ALTER TABLE orders ADD COLUMN shipping_status shipping_status not null default 'unshipped';
ALTER TABLE orders ADD COLUMN pay_status pay_status not null default 'unpaid';
ALTER TABLE orders ADD COLUMN order_status order_status not null default 'unconfirmed';
DROP type order_status;
CREATE type order_status AS ENUM ('unconfirmed', 'confirmed', 'cancelled', 'invalid', 'returned');
