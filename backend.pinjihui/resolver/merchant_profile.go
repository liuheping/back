package resolver

import (
	"errors"
	"strings"

	"golang.org/x/net/context"
	gcontext "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

type merchantProfileResolver struct {
	m *model.MerchantProfile
}

func (r *merchantProfileResolver) UserId() string {
	return r.m.UserId
}

func (r *merchantProfileResolver) SocialId() *string {
	return r.m.SocialId
}

func (r *merchantProfileResolver) RepName() *string {
	return r.m.RepName
}

func (r *merchantProfileResolver) CompanyName() *string {
	return r.m.CompanyName
}

func (r *merchantProfileResolver) CompanyAddress() *addressResolver {
	res := addressResolver{r.m.CompanyAddress}
	return &res
}

func (r *merchantProfileResolver) DeliveryAddress() *addressResolver {
	res := addressResolver{r.m.DeliveryAddress}
	return &res
}

func (r *merchantProfileResolver) LicenseImage(ctx context.Context) *string {
	if r.m.LicenseImage == nil {
		return nil
	}
	res := ctx.Value("config").(*gcontext.Config).QiniuCndHost + *r.m.LicenseImage
	return &res
}

func (r *merchantProfileResolver) CompanyImage(ctx context.Context) *[]string {
	if r.m.CompanyImage == nil {
		return nil
	}
	ra := util.Substr(*r.m.CompanyImage)
	res := strings.Split(ra, ",")
	for x, y := range res {
		res[x] = ctx.Value("config").(*gcontext.Config).QiniuCndHost + y
	}
	return &res
}

func (r *merchantProfileResolver) TakeCash(ctx context.Context) (*[]*takeCashResolver, error) {
	if isAuthorized := ctx.Value("is_authorized").(bool); !isAuthorized {
		return nil, errors.New(gcontext.CredentialsError)
	}
	userId := ctx.Value("user_id").(*string)
	takecashprofile, err := ctx.Value("userRepository").(*repository.UserRepository).TakeCashList(*userId)
	if err != nil {
		return nil, err
	}
	l := make([]*takeCashResolver, len(*takecashprofile))
	for i := range l {
		l[i] = &takeCashResolver{(*takecashprofile)[i]}
	}
	return &l, nil
}

func (r *merchantProfileResolver) ShippingAddress(ctx context.Context) (*[]*shippingAddressResolver, error) {
	addr, err := rp.L("address").(*rp.AddressRepository).FindAll(ctx, r.m.UserId)
	if err != nil {
		return nil, err
	}
	l := make([]*shippingAddressResolver, len(addr))
	for i := range l {
		l[i] = &shippingAddressResolver{(addr)[i]}
	}
	return &l, nil
}

func (r *merchantProfileResolver) Lat() *string {
	return r.m.Lat
}

func (r *merchantProfileResolver) Lng() *string {
	return r.m.Lng
}

func (r *merchantProfileResolver) Balance() float64 {
	return r.m.Balance
}

func (r *merchantProfileResolver) Logo(ctx context.Context) *string {
	if r.m.Logo == nil {
		return nil
	}
	res := ctx.Value("config").(*gcontext.Config).QiniuCndHost + *r.m.Logo
	return &res
}

func (r *merchantProfileResolver) Telephone() *string {
	return r.m.Telephone
}

// func (r *merchantProfileResolver) Waiters(ctx context.Context) *[]string {
// 	if r.m.Waiters == nil {
// 		return nil
// 	}
// 	ra := util.Substr(*r.m.Waiters)
// 	res := strings.Split(ra, ",")
// 	return &res
// }
