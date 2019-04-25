package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CreateWaiter(ctx context.Context, args struct {
	NewWa *model.Waiter
}) (*waiterResolver, error) {
	wa, err := rp.L("waiter").(*rp.WaiterRepository).SaveWaiter(ctx, args.NewWa)
	if err != nil {
		return nil, err
	}
	return &waiterResolver{wa}, nil
}

func (r *Resolver) UpdateWaiter(ctx context.Context, args struct {
	ID    graphql.ID
	NewWa *model.Waiter
}) (*waiterResolver, error) {
	args.NewWa.ID = string(args.ID)
	wa, err := rp.L("waiter").(*rp.WaiterRepository).UpdateWaiter(ctx, args.NewWa)
	if err != nil {
		return nil, err
	}
	return &waiterResolver{wa}, nil
}

func (r *Resolver) DeleteWaiter(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("waiter").(*rp.WaiterRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Resolver) CheckWaiter(ctx context.Context, args struct {
	ID        graphql.ID
	Waiter_id graphql.ID
	Remark    *string
}) (bool, error) {
	_, err := rp.L("waiter").(*rp.WaiterRepository).CheckWaiter(ctx, string(args.ID), string(args.Waiter_id), args.Remark)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Resolver) SaydeleteFromZhiMa(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("waiter").(*rp.WaiterRepository).SaydeleteFromZhiMa(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
