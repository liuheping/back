-- +goose Up
COMMENT ON COLUMN orders.address IS '收货地址';
COMMENT ON COLUMN orders.note IS '管理员操作备注';
COMMENT ON COLUMN orders.inv_taxpayer_id IS '纳税人识别号';
COMMENT ON COLUMN orders.inv_url IS '电子发票地址';
COMMENT ON COLUMN orders.used_coupon IS '使用的优惠券';
COMMENT ON COLUMN orders.offer_amount IS '优惠金额';


comment on column user_coupons.limit_amount
is '订单最低限额';
comment on column coupons.type
is '优惠券类型, 跟领取逻辑相关, 暂时有邀请码领券, 无条件领券';
