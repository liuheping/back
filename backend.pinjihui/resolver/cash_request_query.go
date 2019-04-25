package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

// 根据ID获取提现申请
func (r *Resolver) CashRequest(ctx context.Context, args struct {
	ID string
}) (*cashRequestResolver, error) {
	request, err := rp.L("cashrequest").(*rp.CashRequestRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &cashRequestResolver{request}, nil
}

// 获取所有提现申请
func (r *Resolver) CashRequests(ctx context.Context) (*[]*cashRequestResolver, error) {
	requests, err := rp.L("cashrequest").(*rp.CashRequestRepository).FindAll(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*cashRequestResolver, len(requests))
	for i := range l {
		l[i] = &cashRequestResolver{(requests)[i]}
	}
	return &l, nil
}
