package resolver

import (
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
)

type shippingMethodResolver struct {
    m *model.ShippingMethod
}

func (s *shippingMethodResolver) ID() graphql.ID {
    return graphql.ID(s.m.ID)
}

func (s *shippingMethodResolver) Name() string {
    return s.m.Name
}
