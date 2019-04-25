package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
)

type adResolver struct {
	m *model.Ad
}

func (r *adResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *adResolver) Image(ctx context.Context) string {
	if r.m.Image == "" {
		return ""
	}
	// res := ctx.Value("config").(*gc.Config).QiniuCndHost + r.m.Image
	// return res
	return CompleteUrl(ctx, r.m.Image)
}

func (r *adResolver) Link() *string {
	return r.m.Link
}

func (r *adResolver) Merchant(ctx context.Context) (*merchantProfileResolver, error) {
	if r.m.Merchant_id == nil {
		return nil, nil
	}
	mp, err := ctx.Value("userRepository").(*repository.UserRepository).MerchantProfile(*r.m.Merchant_id)
	if err != nil {
		return nil, err
	}
	return &merchantProfileResolver{mp}, nil
}

func (r *adResolver) Sort() int32 {
	return r.m.Sort
}

func (r *adResolver) Position() string {
	return r.m.Position
}

func (r *adResolver) Isshow() bool {
	return r.m.Is_show
}
