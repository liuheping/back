package resolver

import (
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/service"
)

type operationLogConnectionResolver struct {
	m          []*model.OperationLog
	totalCount int
	from       *string
	to         *string
	hasNext    bool
}

func (r *operationLogConnectionResolver) TotalCount() int32 {
	return int32(r.totalCount)
}

func (r *operationLogConnectionResolver) Edges() *[]*operationLogEdgeResolver {
	res := make([]*operationLogEdgeResolver, len(r.m))
	for i := range res {
		res[i] = &operationLogEdgeResolver{service.EncodeCursor(&r.m[i].ID), r.m[i]}
	}
	return &res
}

func (r *operationLogConnectionResolver) OperationLogs() *[]*operationLogResolver {
	res := make([]*operationLogResolver, len(r.m))
	for i := range res {
		res[i] = &operationLogResolver{r.m[i]}
	}
	return &res
}

func (r *operationLogConnectionResolver) PageInfo() *pageInfoResolver {
	res := pageInfoResolver{
		startCursor: service.EncodeCursor(r.from),
		endCursor:   service.EncodeCursor(r.to),
		hasNextPage: r.hasNext}
	return &res
}
