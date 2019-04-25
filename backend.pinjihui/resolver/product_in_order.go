package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

type productInOrderResolver struct {
	m *model.ProductInOrder
}

func (r *productInOrderResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *productInOrderResolver) Order(ctx context.Context) (*orderResolver, error) {
	order, err := rp.L("order").(*rp.OrderRepository).FindByID(ctx, r.m.Order_id)
	if err != nil {
		return nil, err
	}
	return &orderResolver{order}, nil
}

func (r *productInOrderResolver) Product() (*productResolver, error) {
	product, err := rp.L("product").(*rp.ProductRepository).FindByID(r.m.Product_id)
	if err != nil {
		return nil, err
	}
	return &productResolver{product}, nil
}

func (r *productInOrderResolver) Name() string {
	return r.m.Product_name
}

func (r *productInOrderResolver) ProductCount() int32 {
	return r.m.Product_number
}

func (r *productInOrderResolver) Price() float64 {
	return r.m.Product_price
}

func (r *productInOrderResolver) Image(ctx context.Context) string {
	if r.m.Product_image == "" {
		return ""
	}
	// res := ctx.Value("config").(*gc.Config).QiniuCndHost + r.m.Product_image
	// return res
	return CompleteUrl(ctx, r.m.Product_image)
}

func (r *productInOrderResolver) Specname1() *string {
	return r.m.Spec_1_name
}

func (r *productInOrderResolver) Specname2() *string {
	return r.m.Spec_2_name
}

func (r *productInOrderResolver) Spec1() *string {
	return r.m.Spec_1
}

func (r *productInOrderResolver) Spec2() *string {
	return r.m.Spec_2
}

func (r *productInOrderResolver) ShippingFee() *float64 {
	return r.m.Shipping_fee
}
