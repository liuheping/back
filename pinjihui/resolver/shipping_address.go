package resolver

import (
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
)

type shippingAddressResolver struct {
    m *model.ShippingAddress
}

func (r *shippingAddressResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *shippingAddressResolver) Mobile() string {
    return r.m.Mobile
}

func (r *shippingAddressResolver) Consignee() string {
    return r.m.Consignee
}

func (r *shippingAddressResolver) Address() (*addressResolver, error) {
    return &addressResolver{r.m.Address}, nil
}

func (r *shippingAddressResolver) Zipcode() *string {
    return r.m.Zipcode
}

func (r *shippingAddressResolver) IsDefault() bool {
    return r.m.IsDefault
}
