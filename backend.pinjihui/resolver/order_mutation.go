package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

// 设置订单状态
func (r *Resolver) SetOrderStatus(ctx context.Context, args struct {
	ID     graphql.ID
	Status string
}) (bool, error) {
	_, err := rp.L("order").(*rp.OrderRepository).SetStatus(ctx, string(args.ID), args.Status)
	if err != nil {
		return false, err
	}
	return true, nil
}
