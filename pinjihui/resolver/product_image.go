
package resolver

import (
    //"pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/util"
    "golang.org/x/net/context"
)

type productImageResolver struct {
    m *model.ProductImage
}

func (r *productImageResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *productImageResolver) SamllImage(ctx context.Context) string {
    return getThumbnail(completeUrl(ctx, util.GetString(r.m.SmallImage, r.m.BigImage)))
}

func (r *productImageResolver) MediumImage(ctx context.Context) string {
    return completeUrl(ctx, util.GetString(r.m.MediumImage, r.m.BigImage)) + "?imageView2/0/w/400/h/400"
}

func (r *productImageResolver) BigImage(ctx context.Context) string {
    return completeUrl(ctx, r.m.BigImage) + "?imageView2/0/w/800/h/800"
}
