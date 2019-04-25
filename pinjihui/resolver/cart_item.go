package resolver

import (
    "pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/loader"
    "golang.org/x/net/context"
	"pinjihui.com/pinjihui/util"
	rp "pinjihui.com/pinjihui/repository"
)

type cartItemResolver struct {
    m *model.CartItem
}

func (r *cartItemResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *cartItemResolver) Product(ctx context.Context) (*productResolver, error) {
    product, err := loader.LoadProduct(ctx, r.m.ProductID, r.m.MerchantID)
    if err != nil {
        return nil, err
    }
    return &productResolver{product}, nil
}

func (r *cartItemResolver) ProductCount() int32 {
    return r.m.ProductCount
}

func (r *cartItemResolver) Merchant(ctx context.Context) (*merchantResolver, error) {
    merchant, err := loader.LoadMerchant(ctx, r.m.MerchantID)

    if err != nil {
        return nil, err
    }
    return &merchantResolver{merchant}, nil
}

func (r *cartItemResolver) TotalPrice(ctx context.Context) (string, error) {
    product, err := loader.LoadProduct(ctx, r.m.ProductID, r.m.MerchantID)
    if err != nil {
        return "", err
    }
    price := float64(r.m.ProductCount) * rp.L("product").(*rp.ProductRepository).GetPrice(product, ctx)
    return util.FmtMoney(price), nil
}
