package resolver

import (
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/service"
)

type commentsConnectionResolver struct {
	m          []*model.Comment
	totalCount int
	from       *string
	to         *string
	hasNext    bool
}

func (r *commentsConnectionResolver) TotalCount() int32 {
	return int32(r.totalCount)
}

func (r *commentsConnectionResolver) Edges() *[]*commentsEdgeResolver {
	res := make([]*commentsEdgeResolver, len(r.m))
	for i := range res {
		res[i] = &commentsEdgeResolver{service.EncodeCursor(&r.m[i].ID), r.m[i]}
	}
	return &res
}

func (r *commentsConnectionResolver) Comments() *[]*commentResolver {
	res := make([]*commentResolver, len(r.m))
	for i := range res {
		res[i] = &commentResolver{r.m[i]}
	}
	return &res
}

func (r *commentsConnectionResolver) PageInfo() *pageInfoResolver {
	res := pageInfoResolver{
		startCursor: service.EncodeCursor(r.from),
		endCursor:   service.EncodeCursor(r.to),
		hasNextPage: r.hasNext}
	return &res
}
