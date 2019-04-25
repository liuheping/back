-- +goose Up

CREATE TABLE "public"."rel_agents_products" (
id VARCHAR(32) NOT NULL CONSTRAINT rel_agents_products_pkey PRIMARY KEY,
agent_id VARCHAR(32) NOT NULL,
product_id VARCHAR(32) NOT NULL,
sales_volume int4 NOT NULL DEFAULT 0,
is_sale bool NOT NULL DEFAULT TRUE,
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT "rel_agents_products_merchant_profiles_user_id_fk" FOREIGN KEY ( "agent_id" ) REFERENCES "public"."merchant_profiles" ( "user_id" ) ON DELETE NO ACTION ON UPDATE NO ACTION,
CONSTRAINT "rel_agents_products_products_id_fk" FOREIGN KEY ( "product_id" ) REFERENCES "public"."products" ( "id" ) ON DELETE NO ACTION ON UPDATE NO ACTION 
);

COMMENT ON TABLE "public"."rel_agents_products" IS '代理商产品关联表';



-- +goose Down
DROP TABLE "public"."rel_agents_products"