
package resolver
import (
    //"pinjihui.com/backend.pinjihui/model"
    "github.com/graph-gophers/graphql-go"
)
type roleResolver struct {
    //m *model.role
}

func (r *roleResolver) ID() graphql.ID {
    res := graphql.ID("xjauwkahsi92h1j")
    return res
}

func (r *roleResolver) Name() *string {
    res := "test string"
    return &res
}
