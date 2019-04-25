package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

func (r *Resolver) Order(ctx context.Context, args struct {
	ID string
}) (*orderResolver, error) {
	con, err := rp.L("order").(*rp.OrderRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &orderResolver{con}, nil
}

//根据条件获取所有订单
func (r *Resolver) Orders(ctx context.Context, args struct {
	First  *int32
	Offset *int32
	Search *model.OrderSearchInput
	Sort   *model.OrderSortInput
}) (*ordersConnectionResolver, error) {
	fetchSize := util.GetInt32(args.First, DefaultPageSize)
	list, err := rp.L("order").(*rp.OrderRepository).Search(ctx, &fetchSize, args.Offset, args.Search, args.Sort)
	if err != nil {
		return nil, err
	}
	count, err := rp.L("order").(*rp.OrderRepository).Count(ctx, args.Search)
	if err != nil {
		return nil, err
	}
	var from *string
	var to *string
	if len(list) > 0 {
		from = &(list[0].ID)
		to = &(list[len(list)-1].ID)
	}
	return &ordersConnectionResolver{m: list, totalCount: count, from: from, to: to, hasNext: util.If(len(list) == int(fetchSize), true, false).(bool)}, nil
}

//根据条件和加盟商所在区域获取所有商品
func (r *Resolver) OrdersByArea(ctx context.Context, args struct {
	First  *int32
	Offset *int32
	Search *model.OrderSearchInput
	Sort   *model.OrderSortInput
}) (*ordersConnectionResolver, error) {
	fetchSize := util.GetInt32(args.First, DefaultPageSize)
	list, err := rp.L("order").(*rp.OrderRepository).SearchByArea(ctx, &fetchSize, args.Offset, args.Search, args.Sort)
	if err != nil {
		return nil, err
	}
	count, err := rp.L("order").(*rp.OrderRepository).CountByArea(ctx, args.Search)
	if err != nil {
		return nil, err
	}
	var from *string
	var to *string
	if len(list) > 0 {
		from = &(list[0].ID)
		to = &(list[len(list)-1].ID)
	}
	return &ordersConnectionResolver{m: list, totalCount: count, from: from, to: to, hasNext: util.If(len(list) == int(fetchSize), true, false).(bool)}, nil
}
