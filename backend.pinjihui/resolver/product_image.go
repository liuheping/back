package resolver

import (
	"time"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
)

type productImageResolver struct {
	m *model.ProductImage
}

func (r *productImageResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *productImageResolver) ProductID() graphql.ID {
	return graphql.ID(r.m.Product_id)
}

func (r *productImageResolver) SamllImage(ctx context.Context) *string {
	if r.m.Big_image == "" {
		return nil
	}
	// res := ctx.Value("config").(*gc.Config).QiniuCndHost + r.m.Big_image + "?imageView2/2/w/160/h/240"
	// return &res
	url := CompleteUrl(ctx, r.m.Big_image) + "?imageView2/2/w/200/h/200"
	return &url
}

func (r *productImageResolver) MediumImage(ctx context.Context) *string {
	if r.m.Big_image == "" {
		return nil
	}
	// res := ctx.Value("config").(*gc.Config).QiniuCndHost + r.m.Big_image + "?imageView2/2/w/320/h/480"
	// return &res
	url := CompleteUrl(ctx, r.m.Big_image) + "?imageView2/2/w/400/h/400"
	return &url
}

func (r *productImageResolver) BigImage(ctx context.Context) string {
	if r.m.Big_image == "" {
		return ""
	}
	// res := ctx.Value("config").(*gc.Config).QiniuCndHost + r.m.Big_image
	// return res
	return CompleteUrl(ctx, r.m.Big_image)
}

func (r *productImageResolver) Created_at() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Created_at)
	return graphql.Time{Time: res}, err
}
