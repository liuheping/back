package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
)

type operationLogEdgeResolver struct {
	cursor graphql.ID
	m      *model.OperationLog
}

func (r *operationLogEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *operationLogEdgeResolver) Node() *operationLogResolver {
	return &operationLogResolver{r.m}
}
