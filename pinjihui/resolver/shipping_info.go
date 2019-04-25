package resolver

import (
    "pinjihui.com/pinjihui/model"
    "golang.org/x/net/context"
)

type shippingInfoResolver struct {
    m *model.ShippingInfo
}

func (s *shippingInfoResolver) Company() *string {
    return s.m.Company
}

func (s *shippingInfoResolver) Number() *string {
    return s.m.DeliveryNumber
}

func (s *shippingInfoResolver) Images(ctx context.Context) *[]string {
    arr := s.m.ImageArr()
    for i, url := range arr {
        arr[i] = completeUrl(ctx, url)
    }
    return &arr
}
