package resolver

import (
    "time"

    "github.com/graph-gophers/graphql-go"
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/model"
    rp "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/service"
)

type userResolver struct {
    u *model.User
}

func (r *userResolver) ID() graphql.ID {
    return graphql.ID(r.u.ID)
}

func (r *userResolver) Name() *string {
    return r.u.Name
}

func (r *userResolver) Mobile() *string {
    return r.u.Mobile
}

func (r *userResolver) Email() *string {
    return r.u.Email
}

func (r *userResolver) CreatedAt() (graphql.Time, error) {
    t, err := time.Parse(time.RFC3339, r.u.CreatedAt)
    return graphql.Time{Time: t}, err
}

func (r *userResolver) LastLoginTime() (*graphql.Time, error) {
    if r.u.LastLoginTime == nil {
        return nil, nil
    }

    t, err := time.Parse(time.RFC3339, *r.u.LastLoginTime)
    return &graphql.Time{Time: t}, err
}

func (r *userResolver) LastIp() *string {
    return &r.u.LastIp
}

//这里如果没有记录,不应该返回错误,而是返回(nil,nil)
func (r *userResolver) Addresses(ctx context.Context) (*[]*shippingAddressResolver, error) {
    addrs, err := rp.L("address").(*rp.AddressRepository).FindAll(ctx)
    if err != nil {
        return nil, err
    }
    addrsr := make([]*shippingAddressResolver, len(addrs))
    for i, v := range addrs {
        addrsr[i] = &shippingAddressResolver{v}
    }
    return &addrsr, nil
}

func (r *userResolver) Cart(ctx context.Context) ([]*cartItemResolver, error) {
    return getCartItemResolver(ctx)
}

func (r *userResolver) Token(ctx context.Context) (*string, error) {
    return ctx.Value("authService").(*service.AuthService).SignJWT(r.u)
}

func (r *userResolver) InviteCode() string {
    return r.u.InviteCode
}

func (r *userResolver) Invited() bool {
    return r.u.Invited
}

func (r *userResolver) IsFirstLogin() bool {
    return r.u.LastLoginTime == nil
}

func (r *userResolver) Type() string {
    return r.u.Type
}

func (r *userResolver) HasShareCoupon(ctx context.Context) (bool, error) {
    return rp.L("user").(*rp.UserRepository).HasShareCoupon(ctx, model.ForSharer)
}

func (r *userResolver) HasBeShareCoupon(ctx context.Context) (bool, error) {
    return rp.L("user").(*rp.UserRepository).HasShareCoupon(ctx, model.ForBeSharer)
}
