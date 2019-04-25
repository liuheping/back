package resolver

import (
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
	rp "pinjihui.com/backend.pinjihui/repository"
)

type spikeResolver struct {
	m *model.Spike
}

func (r *spikeResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *spikeResolver) Price() float64 {
	return r.m.Price
}

func (r *spikeResolver) StartAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Start_at)
	return graphql.Time{Time: res}, err
}

func (r *spikeResolver) ExpiredAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Expired_at)
	return graphql.Time{Time: res}, err
}

func (r *spikeResolver) TotalCount() int32 {
	return r.m.Total_count
}

func (r *spikeResolver) BuyLimit() int32 {
	return r.m.Buy_limit
}

func (r *spikeResolver) Merchant(ctx context.Context) (*merchantProfileResolver, error) {
	mp, err := ctx.Value("userRepository").(*repository.UserRepository).MerchantProfile(r.m.Merchant_id)
	if err != nil {
		return nil, err
	}
	return &merchantProfileResolver{mp}, nil
}

func (r *spikeResolver) Product() (*productResolver, error) {
	product, err := rp.L("product").(*rp.ProductRepository).FindByID(r.m.Product_id)
	if err != nil {
		return nil, err
	}
	return &productResolver{product}, nil
}
