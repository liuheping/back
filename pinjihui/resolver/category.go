
package resolver

import (
    "pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
    "golang.org/x/net/context"
)

type categoryResolver struct {
    m *model.Category
}

func (r *categoryResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *categoryResolver) Name() string {
    return r.m.Name
}

func (r *categoryResolver) ParentId() *graphql.ID {
    if r.m.ParentId == nil {
        return nil
    }
    p := graphql.ID(*r.m.ParentId)
    return &p
}

func (r *categoryResolver) Thumbnail(ctx context.Context) *string {
    if r.m.Thumbnail != nil {
        url := completeUrl(ctx, *r.m.Thumbnail)
        return &url
    }
    return nil
}

func (r *categoryResolver) Children() ([]*categoryResolver, error) {
    return GetChildResolvers(r.m)
}

func GetChildResolvers(parent *model.Category) ([]*categoryResolver, error) {
    listr := make([]*categoryResolver, len(parent.Children))
    for i, v := range parent.Children {
        listr[i] = &categoryResolver{v}
    }
    return listr, nil
}

func (r *categoryResolver) IsCommon() bool {
    return r.m.IsCommon
}
