package resolver

import (
    "pinjihui.com/pinjihui/model"
)

type productSpecResolver struct {
    p *model.PaMCPair
    items []*model.ProductSpec
}

func (r *productSpecResolver) SpecItems() ([]*productSpecItemResolver, error) {
    pr := make([]*productSpecItemResolver, len(r.items))
    for i, v := range r.items {
        pr[i] = &productSpecItemResolver{v}
    }
    return pr, nil
}

func (r *productSpecResolver) Spec1Name() string {
    return *r.p.Spec1Name
}

func (r *productSpecResolver) Spec2Name() *string {
    return r.p.Spec2Name
}
