-- +goose Up
ALTER TABLE waiters ADD COLUMN remark VARCHAR(500);

comment on table waiters is '客服表';
comment on column waiters.name is '客服姓名';
comment on column waiters.mobile is '客服电话';
comment on column waiters.waiter_id is '芝麻客服ID';
comment on column waiters.handled is '商家删除客服之后，标记管理员是否删除芝麻客服账号（只有当deleted为true的时候才有意义）';


CREATE UNIQUE INDEX waiters_mobile_uindex ON public.waiters (mobile);
ALTER TABLE public.waiters ADD CONSTRAINT waiters_merchant_profiles_user_id_fk FOREIGN KEY (merchant_id) REFERENCES merchant_profiles (user_id);

-- +goose Down
ALTER TABLE waiters DROP COLUMN remark;
DROP INDEX users_invite_code_uindex;
DROP CONSTRAINT waiters_merchant_profiles_user_id_fk;