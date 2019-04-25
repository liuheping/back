package resolver

import (
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	gcontext "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

type shippingInfoResolver struct {
	m *model.ShippingInfo
}

func (r *shippingInfoResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *shippingInfoResolver) Company() *string {
	return r.m.Company
}

func (r *shippingInfoResolver) Delivery_number() *string {
	return r.m.Delivery_number
}

func (r *shippingInfoResolver) Images(ctx context.Context) *[]string {
	if r.m.Images == nil {
		return nil
	}
	ra := util.Substr(*r.m.Images)
	res := strings.Split(ra, ",")
	for x, y := range res {
		res[x] = ctx.Value("config").(*gcontext.Config).QiniuCndHost + y
	}
	return &res
}

func (r *shippingInfoResolver) Order(ctx context.Context) (*orderResolver, error) {
	con, err := rp.L("order").(*rp.OrderRepository).FindByID(ctx, r.m.Order_id)
	if err != nil {
		return nil, err
	}
	return &orderResolver{con}, nil
}
