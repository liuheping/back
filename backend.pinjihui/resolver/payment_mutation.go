package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CreatePayment(ctx context.Context, args struct {
	NewPay *model.Payment
}) (*paymentResolver, error) {
	pay, err := rp.L("payment").(*rp.PaymentRepository).SavePayment(ctx, args.NewPay)
	if err != nil {
		return nil, err
	}
	return &paymentResolver{pay}, nil
}

func (r *Resolver) UpdatePayment(ctx context.Context, args struct {
	ID     graphql.ID
	NewPay *model.Payment
}) (*paymentResolver, error) {
	args.NewPay.ID = string(args.ID)
	pay, err := rp.L("payment").(*rp.PaymentRepository).SavePayment(ctx, args.NewPay)
	if err != nil {
		return nil, err
	}
	return &paymentResolver{pay}, nil
}

func (r *Resolver) DeletePayment(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("payment").(*rp.PaymentRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
