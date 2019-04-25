package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CreateCoupon(ctx context.Context, args struct {
	Cou *model.Coupon
}) (*couponResolver, error) {
	con, err := rp.L("coupon").(*rp.CouponRepository).SaveCoupon(ctx, args.Cou)
	if err != nil {
		return nil, err
	}
	return &couponResolver{con}, nil
}

func (r *Resolver) UpdateCoupon(ctx context.Context, args struct {
	ID  graphql.ID
	Cou *model.Coupon
}) (*couponResolver, error) {
	args.Cou.ID = string(args.ID)
	con, err := rp.L("coupon").(*rp.CouponRepository).SaveCoupon(ctx, args.Cou)
	if err != nil {
		return nil, err
	}
	return &couponResolver{con}, nil
}

func (r *Resolver) DeleteCoupon(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("coupon").(*rp.CouponRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
