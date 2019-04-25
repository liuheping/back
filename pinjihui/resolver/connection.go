package resolver

import "pinjihui.com/pinjihui/service"

type Connection struct {
    totalCount int
    from       *string
    to         *string
    hasNext    bool
}

func (r *Connection) PageInfo() *pageInfoResolver {
    res := pageInfoResolver{
        startCursor: service.EncodeCursor(r.from),
        endCursor:   service.EncodeCursor(r.to),
        hasNextPage: r.hasNext}
    return &res
}


func (r *Connection) TotalCount() int32 {
    return int32(r.totalCount)
}
