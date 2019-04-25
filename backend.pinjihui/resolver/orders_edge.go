package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
)

type ordersEdgeResolver struct {
	cursor graphql.ID
	m      *model.Order
}

func (r *ordersEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *ordersEdgeResolver) Node() *orderResolver {
	return &orderResolver{r.m}
}
