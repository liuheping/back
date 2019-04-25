package resolver

import (
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/service"
)

type waiterConnectionResolver struct {
	m          []*model.Waiter
	totalCount int
	from       *string
	to         *string
	hasNext    bool
}

func (r *waiterConnectionResolver) TotalCount() int32 {
	return int32(r.totalCount)
}

func (r *waiterConnectionResolver) Edges() *[]*waiterEdgeResolver {
	res := make([]*waiterEdgeResolver, len(r.m))
	for i := range res {
		res[i] = &waiterEdgeResolver{service.EncodeCursor(&r.m[i].ID), r.m[i]}
	}
	return &res
}

func (r *waiterConnectionResolver) Waiters() *[]*waiterResolver {
	res := make([]*waiterResolver, len(r.m))
	for i := range res {
		res[i] = &waiterResolver{r.m[i]}
	}
	return &res
}

func (r *waiterConnectionResolver) PageInfo() *pageInfoResolver {
	res := pageInfoResolver{
		startCursor: service.EncodeCursor(r.from),
		endCursor:   service.EncodeCursor(r.to),
		hasNextPage: r.hasNext}
	return &res
}
