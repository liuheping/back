package resolver

import (
    //"pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/model"
)

type commentsConnectionResolver struct {
    list []*model.Comment
    product *model.Product
    Connection
}

func (r *commentsConnectionResolver) Comments() *[]*commentResolver {
    res := make([]*commentResolver, len(r.list))
    for i := range res {
        res[i] = &commentResolver{r.list[i], r.product}
    }
    return &res
}
