
package resolver
import (
    //"pinjihui.com/backend.pinjihui/model"
    "github.com/graph-gophers/graphql-go"
)
type attributesEdgeResolver struct {
    //m *model.attributesEdge
}

func (r *attributesEdgeResolver) Cursor() graphql.ID {
    res := graphql.ID("xjauwkahsi92h1j")
    return res
}

func (r *attributesEdgeResolver) Node() *attributeResolver {
    res := attributeResolver{}
    return &res
}
