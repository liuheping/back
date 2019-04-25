-- +goose Up
ALTER TABLE carts ADD COLUMN agent_id varchar(32) NULL;
ALTER TABLE order_products ADD COLUMN agent_id varchar(32) NULL;
CREATE UNIQUE INDEX rel_agents_products_agent_id_product_id_uindex ON rel_agents_products (agent_id, product_id);
-- +goose Down
ALTER TABLE carts DROP COLUMN agent_id;
ALTER TABLE order_products DROP COLUMN agent_id;
DROP INDEX rel_agents_products_agent_id_product_id_uindex;
