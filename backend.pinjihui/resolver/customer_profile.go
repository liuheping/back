package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

type customerProfileResolver struct {
	m *model.WechartProfile
}

func (r *customerProfileResolver) WechartProfile(ctx context.Context) (*wechartProfileResolver, error) {
	if r.m == nil {
		return nil, nil
	}
	profile, err := rp.L("wechartprofile").(*rp.WechartProfileRepository).FindByUserID(ctx, r.m.User_id)
	if err != nil {
		return nil, err
	}
	if profile == nil && err == nil {
		return nil, nil
	}
	return &wechartProfileResolver{profile}, nil
}

func (r *customerProfileResolver) ShippingAddress(ctx context.Context) (*[]*shippingAddressResolver, error) {
	addr, err := rp.L("address").(*rp.AddressRepository).FindAll(ctx, r.m.User_id)
	if err != nil {
		return nil, err
	}
	l := make([]*shippingAddressResolver, len(addr))
	for i := range l {
		l[i] = &shippingAddressResolver{(addr)[i]}
	}
	return &l, nil
}
