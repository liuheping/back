package resolver

import (
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/service"
)

type usersConnectionResolver struct {
	m          []*model.User
	totalCount int
	from       *string
	to         *string
	hasNext    bool
}

func (r *usersConnectionResolver) TotalCount() int32 {
	return int32(r.totalCount)
}

func (r *usersConnectionResolver) Edges() *[]*usersEdgeResolver {
	res := make([]*usersEdgeResolver, len(r.m))
	for i := range res {
		res[i] = &usersEdgeResolver{service.EncodeCursor(&r.m[i].ID), r.m[i]}
	}
	return &res
}

func (r *usersConnectionResolver) Users() *[]*userResolver {
	res := make([]*userResolver, len(r.m))
	for i := range res {
		res[i] = &userResolver{r.m[i]}
	}
	return &res
}

func (r *usersConnectionResolver) PageInfo() *pageInfoResolver {
	res := pageInfoResolver{
		startCursor: service.EncodeCursor(r.from),
		endCursor:   service.EncodeCursor(r.to),
		hasNextPage: r.hasNext}
	return &res
}
