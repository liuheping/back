package resolver

import (
    //"pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/service"
)

const DefaultPageSize  = 10

type productsConnectionResolver struct {
    m          []*model.PaMCPair
    Connection
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
