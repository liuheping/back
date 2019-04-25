package resolver

import (
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/service"
)

const DefaultPageSize = 10

type productsConnectionResolver struct {
	m          []*model.Product
	totalCount int
	from       *string
	to         *string
	hasNext    bool
}

func (r *productsConnectionResolver) TotalCount() int32 {
	return int32(r.totalCount)
}

func (r *productsConnectionResolver) Edges() *[]*productsEdgeResolver {
	res := make([]*productsEdgeResolver, len(r.m))
	for i := range res {
		res[i] = &productsEdgeResolver{service.EncodeCursor(&r.m[i].ID), r.m[i]}
	}
	return &res
}

func (r *productsConnectionResolver) Products() *[]*productResolver {
	res := make([]*productResolver, len(r.m))
	for i := range res {
		res[i] = &productResolver{r.m[i]}
	}
	return &res
}

func (r *productsConnectionResolver) PageInfo() *pageInfoResolver {
	res := pageInfoResolver{
		startCursor: service.EncodeCursor(r.from),
		endCursor:   service.EncodeCursor(r.to),
		hasNextPage: r.hasNext}
	return &res
}
