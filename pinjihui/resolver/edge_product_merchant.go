package resolver

import (
    "pinjihui.com/pinjihui/model"
)

type edgeProductMerchantResolver struct {
    m *model.MerchantWithStock

}

func (r *edgeProductMerchantResolver) Merchant() *merchantResolver {
    return &merchantResolver{&r.m.Merchant}
}

func (r *edgeProductMerchantResolver) Stock() int32 {
    return r.m.Stock
}

func (r *edgeProductMerchantResolver) Price() float64 {
    return r.m.Price
}

func (r *edgeProductMerchantResolver) Distance() string {
    return FmtDistance(r.m.Distance)
}
