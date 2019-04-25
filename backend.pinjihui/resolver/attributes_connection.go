
package resolver
import (
    //"pinjihui.com/backend.pinjihui/model"
)
type attributesConnectionResolver struct {
    //m *model.attributesConnection
}

func (r *attributesConnectionResolver) TotalCount() int32 {
    res := int32(3)
    return res
}

func (r *attributesConnectionResolver) Edges() *[]*attributesEdgeResolver {
    res := make([]*attributesEdgeResolver, 3)
    for i := range res {
        v := attributesEdgeResolver{}
        res[i] = &v
    }
    return &res
}

func (r *attributesConnectionResolver) PageInfo() *pageInfoResolver {
    res := pageInfoResolver{}
    return &res
}
