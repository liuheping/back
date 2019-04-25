package resolver

import (
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/service"
)

type ordersConnectionResolver struct {
	m          []*model.Order
	totalCount int
	from       *string
	to         *string
	hasNext    bool
}

func (r *ordersConnectionResolver) TotalCount() int32 {
	return int32(r.totalCount)
}

func (r *ordersConnectionResolver) Edges() *[]*ordersEdgeResolver {
	res := make([]*ordersEdgeResolver, len(r.m))
	for i := range res {
		res[i] = &ordersEdgeResolver{service.EncodeCursor(&r.m[i].ID), r.m[i]}
	}
	return &res
}

func (r *ordersConnectionResolver) Orders() *[]*orderResolver {
	res := make([]*orderResolver, len(r.m))
	for i := range res {
		res[i] = &orderResolver{r.m[i]}
	}
	return &res
}

func (r *ordersConnectionResolver) PageInfo() *pageInfoResolver {
	res := pageInfoResolver{
		startCursor: service.EncodeCursor(r.from),
		endCursor:   service.EncodeCursor(r.to),
		hasNextPage: r.hasNext}
	return &res
}
