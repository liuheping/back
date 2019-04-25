
package resolver

import (
    //"pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
)

type productsEdgeResolver struct {
    cursor graphql.ID
    m *model.PaMCPair
}

func (r *productsEdgeResolver) Cursor() graphql.ID {
    return r.cursor
}

func (r *productsEdgeResolver) Node() *productResolver {
    return &productResolver{r.m}
}
