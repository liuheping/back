package resolver

import (
	"time"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
)

type couponResolver struct {
	m *model.Coupon
}

func (r *couponResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *couponResolver) Description() string {
	return r.m.Description
}

func (r *couponResolver) Value() float64 {
	return r.m.Value
}

func (r *couponResolver) CreatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Created_at)
	return graphql.Time{Time: res}, err
}

func (r *couponResolver) UpdatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Updated_at)
	return graphql.Time{Time: res}, err
}

func (r *couponResolver) LimitAmount() *float64 {
	return r.m.Limit_amount
}

func (r *couponResolver) ExpiredAt() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, *r.m.Expired_at)
	return &graphql.Time{Time: res}, err
}

func (r *couponResolver) Quantity() int32 {
	return r.m.Quantity
}

func (r *couponResolver) Type() string {
	return r.m.Type
}

func (r *couponResolver) StartAt() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, *r.m.Start_at)
	return &graphql.Time{Time: res}, err
}

func (r *couponResolver) Merchant(ctx context.Context) (*merchantProfileResolver, error) {
	if r.m.Merchant_id == nil {
		return nil, nil
	}
	mp, err := ctx.Value("userRepository").(*repository.UserRepository).MerchantProfile(*r.m.Merchant_id)
	if err != nil {
		return nil, err
	}
	return &merchantProfileResolver{mp}, nil
}

func (r *couponResolver) Validity_days() *int32 {
	return r.m.Validity_days
}
