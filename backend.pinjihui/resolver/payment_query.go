package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) Payment(ctx context.Context, args struct {
	ID string
}) (*paymentResolver, error) {
	pay, err := rp.L("payment").(*rp.PaymentRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &paymentResolver{pay}, nil
}

// //获取所有品牌
func (r *Resolver) Payments(ctx context.Context) (*[]*paymentResolver, error) {
	pay, err := rp.L("payment").(*rp.PaymentRepository).FindAll(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*paymentResolver, len(pay))
	for i := range l {
		l[i] = &paymentResolver{(pay)[i]}
	}
	return &l, nil
}
