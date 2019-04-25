package resolver

import (
    "pinjihui.com/pinjihui/model"
)

type ordersConnectionResolver struct {
    list []*model.Order
    Connection
}

/*func (r *ordersConnectionResolver) Edges() *[]*ordersEdgeResolver {
    res := make([]*ordersEdgeResolver, 3)
    for i := range res {
        v := ordersEdgeResolver{}
        res[i] = &v
    }
    return &res
}*/

func (r *ordersConnectionResolver) Orders() (*[]*orderResolver, error) {

    res := make([]*orderResolver, len(r.list))
    for i := range res {
        res[i] = &orderResolver{m: r.list[i]}
    }
    return &res, nil
}
