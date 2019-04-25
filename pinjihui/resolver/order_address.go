
package resolver

import (

    "pinjihui.com/pinjihui/model"
)

type orderAddressResolver struct {
    m *model.OrderAddress
}

func (r *orderAddressResolver) Mobile() string {
    return r.m.Mobile
}

func (r *orderAddressResolver) Consignee() string {
    return r.m.Consignee
}

func (r *orderAddressResolver) Address() string {
    return *r.m.RegionName + r.m.Address
}

func (r *orderAddressResolver) Zipcode() *string {
    return r.m.Zipcode
}
