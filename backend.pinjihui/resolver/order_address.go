package resolver

import (
	"pinjihui.com/backend.pinjihui/model"
)

type orderAddressResolver struct {
	m *model.OrderAddress
}

func (r *orderAddressResolver) RegionName() *string {
	return r.m.RegionName
}

func (r *orderAddressResolver) Mobile() string {
	return r.m.Mobile
}

func (r *orderAddressResolver) Consignee() string {
	return r.m.Consignee
}

func (r *orderAddressResolver) Address() string {
	return r.m.Address
}

func (r *orderAddressResolver) Zipcode() *string {
	return r.m.Zipcode
}

func (r *orderAddressResolver) AreaID() int32 {
	return r.m.AreaID
}
