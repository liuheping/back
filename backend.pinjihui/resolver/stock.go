package resolver

import (
	"time"

	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
)

type stockResolver struct {
	s *model.Stock
}

func (r *stockResolver) ProductId() graphql.ID {
	return graphql.ID(r.s.Product_id)
}

func (r *stockResolver) MerchantId() graphql.ID {
	return graphql.ID(r.s.Merchant_id)
}

func (r *stockResolver) Stock() int32 {
	return r.s.Stock
}

func (r *stockResolver) RetailPrice() float64 {
	return r.s.Retail_price
}

func (r *stockResolver) SalesVolume() int32 {
	return r.s.Sales_volume
}

func (r *stockResolver) CreatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.s.Created_at)
	return graphql.Time{Time: res}, err
}
