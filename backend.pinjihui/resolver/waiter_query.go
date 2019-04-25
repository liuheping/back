package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

//根据条件获取所有客服
func (r *Resolver) Waiters(ctx context.Context, args struct {
	First  *int32
	Offset *int32
	Search *model.WaiterSearchInput
	Sort   *model.WaiterSortInput
}) (*waiterConnectionResolver, error) {
	fetchSize := util.GetInt32(args.First, DefaultPageSize)
	list, err := rp.L("waiter").(*rp.WaiterRepository).Search(ctx, &fetchSize, args.Offset, args.Search, args.Sort)
	if err != nil {
		return nil, err
	}
	count, err := rp.L("waiter").(*rp.WaiterRepository).Count(ctx, args.Search)
	if err != nil {
		return nil, err
	}
	var from *string
	var to *string
	if len(list) > 0 {
		from = &(list[0].ID)
		to = &(list[len(list)-1].ID)
	}
	return &waiterConnectionResolver{m: list, totalCount: count, from: from, to: to, hasNext: util.If(len(list) == int(fetchSize), true, false).(bool)}, nil
}

//根据ID查找客服
func (r *Resolver) Waiter(ctx context.Context, args struct {
	ID string
}) (*waiterResolver, error) {
	wa, err := rp.L("waiter").(*rp.WaiterRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &waiterResolver{wa}, nil
}

// 通过商家ID查找客服
func (r *Resolver) FindWaitersByMerchantID(ctx context.Context) (*[]*waiterResolver, error) {
	waiters, err := rp.L("waiter").(*rp.WaiterRepository).FindByMerchantID(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*waiterResolver, len(waiters))
	for i := range l {
		l[i] = &waiterResolver{(waiters)[i]}
	}
	return &l, nil
}
