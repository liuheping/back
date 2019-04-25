package resolver

import (
    //"pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/service"
)

type merchantsConnectionResolver struct {
    m          []*model.Merchant
    totalCount int
    from       *string
    to         *string
    hasNext    bool
}

func (r *merchantsConnectionResolver) TotalCount() int32 {
    return int32(r.totalCount)
}

func (r *merchantsConnectionResolver) Merchants() []*merchantResolver {
    res := make([]*merchantResolver, len(r.m))
    for i := range res {
        res[i] = &merchantResolver{r.m[i]}
    }
    return res
}

func (r *merchantsConnectionResolver) PageInfo() *pageInfoResolver {
    res := pageInfoResolver{
        startCursor: service.EncodeCursor(r.from),
        endCursor:   service.EncodeCursor(r.to),
        hasNextPage: r.hasNext}
    return &res
}
