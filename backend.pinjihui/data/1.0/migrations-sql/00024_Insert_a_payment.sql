-- +goose Up
INSERT INTO public.payments (id, pay_name, pay_code, pay_fee, pay_desc, sort_order, enabled, is_cod, is_online, deleted) VALUES ('bc5mn7n2oau0r8ddf3r0', '微信', 'WechartPay', '未知', '腾讯集团，最多人的选择', 255, true, false, true, false) ON CONFLICT DO NOTHING;
-- +goose Down
DELETE FROM payments WHERE id='bc5mn7n2oau0r8ddf3r0';