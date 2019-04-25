
package resolver

import (
    "pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
)

type productSpecItemResolver struct {
    m *model.ProductSpec
}

func (r *productSpecItemResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (p *productSpecItemResolver) Spec1() string {
    return p.m.Spec1
}

func (p *productSpecItemResolver) Spec2() *string {
    return p.m.Spec2
}

func (p *productSpecItemResolver) Price() float64 {
    return p.m.Price
}

func (p *productSpecItemResolver) Stock() int32 {
    return p.m.Stock
}
