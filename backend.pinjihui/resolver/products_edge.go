package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
)

type productsEdgeResolver struct {
	cursor graphql.ID
	m      *model.Product
}

func (r *productsEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *productsEdgeResolver) Node() *productResolver {
	return &productResolver{r.m}
}
