package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) Coupon(ctx context.Context, args struct {
	ID string
}) (*couponResolver, error) {
	con, err := rp.L("coupon").(*rp.CouponRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &couponResolver{con}, nil
}

// 获取所有优惠券
func (r *Resolver) Coupons(ctx context.Context) (*[]*couponResolver, error) {
	cons, err := rp.L("coupon").(*rp.CouponRepository).FindAll(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*couponResolver, len(cons))
	for i := range l {
		l[i] = &couponResolver{(cons)[i]}
	}
	return &l, nil
}
