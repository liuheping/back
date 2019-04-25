-- +goose Up
CREATE TYPE inout_type AS ENUM ('order', 'withdrawal', 'commission');
CREATE TABLE merchant_balance_logs (
  id           SERIAL          NOT NULL CONSTRAINT merchant_balance_logs_pkey PRIMARY KEY,
  merchant_id  VARCHAR(32)     NOT NULL,
  inout        NUMERIC(100, 2) NOT NULL,
  created_at   TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
  inout_type   inout_type      NOT NULL,
  "references" VARCHAR(32)     NOT NULL
);
COMMENT ON TABLE public.merchant_balance_logs IS '商户流水日志';
COMMENT ON COLUMN public.merchant_balance_logs.inout IS '收入或支出';
COMMENT ON COLUMN public.merchant_balance_logs.inout_type IS '收入或支出类型,目前有订单收入,订单提成和提现支出';
COMMENT ON COLUMN public.merchant_balance_logs."references" IS '关联的订单或提现申请id';
-- +goose Down
DROP TABLE merchant_balance_logs;
DROP TYPE inout_type;
