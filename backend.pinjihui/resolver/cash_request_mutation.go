package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

// 商家申请提现
func (r *Resolver) ApplyCash(ctx context.Context, args struct {
	Info *model.CashRequest
}) (*cashRequestResolver, error) {
	request, err := rp.L("cashrequest").(*rp.CashRequestRepository).SaveCashRequest(ctx, args.Info)
	if err != nil {
		return nil, err
	}
	return &cashRequestResolver{request}, nil
}

// 管理员拒绝提现
func (r *Resolver) RefusedCash(ctx context.Context, args struct {
	ID    string
	Reply string
}) (bool, error) {
	_, err := rp.L("cashrequest").(*rp.CashRequestRepository).SetRefused(ctx, args.ID, args.Reply)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 管理员设置提现状态为打款
func (r *Resolver) PaidCash(ctx context.Context, args struct {
	ID    string
	Reply string
}) (bool, error) {
	_, err := rp.L("cashrequest").(*rp.CashRequestRepository).SetPaid(ctx, args.ID, args.Reply)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 商家设置提现状态为完成（商家确认）
func (r *Resolver) FinishedCash(ctx context.Context, args struct {
	ID string
}) (bool, error) {
	_, err := rp.L("cashrequest").(*rp.CashRequestRepository).SetFinished(ctx, args.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}
