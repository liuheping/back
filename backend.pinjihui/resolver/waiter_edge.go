package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
)

type waiterEdgeResolver struct {
	cursor graphql.ID
	m      *model.Waiter
}

func (r *waiterEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *waiterEdgeResolver) Node() *waiterResolver {
	return &waiterResolver{r.m}
}
