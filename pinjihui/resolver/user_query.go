package resolver

import (
    rp "pinjihui.com/pinjihui/repository"
    "golang.org/x/net/context"
    gc "pinjihui.com/pinjihui/context"
    "pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/util"
    "pinjihui.com/pinjihui/service"
)

func (r *Resolver) Me(ctx context.Context) (*userResolver, error) {
    gc.CheckAuth(ctx)
    user := gc.User(ctx)
    return &userResolver{user}, nil
}

func (r *Resolver) Login(ctx context.Context, args struct {
    Mobile   string
    Password string
    Code     *string
}) (*userResolver, error) {
    userCredentials := &model.UserCredentials{
        Mobile:   args.Mobile,
        Password: args.Password,
    }
    user, err := rp.L("user").(*rp.UserRepository).ComparePassword(userCredentials)
    if err != nil {
        return nil, err
    }
    if args.Code != nil {
        ctx = context.WithValue(ctx, "user", user)
        ctx = context.WithValue(ctx, "is_authorized", true)
        wxUser, err := rp.L("user").(*rp.UserRepository).WxLogin(ctx, *args.Code)
        if err != nil {
            return nil, err
        }
        user.Openid = wxUser.OpenId
    }
    rp.L("user").(*rp.UserRepository).UpdateLoginStatus(ctx, user.ID)
    return &userResolver{user}, nil
}

func (r *Resolver) MyAddress(ctx context.Context, args struct {
    ID graphql.ID
}) (*shippingAddressResolver, error) {
    ID := string(args.ID)
    addr, err := rp.L("address").(*rp.AddressRepository).FindByID(ctx, ID)
    if err != nil {
        return nil, err
    }
    return &shippingAddressResolver{addr}, nil
}

func (r *Resolver) WxLogin(ctx context.Context, args struct{ Code string }) (*userResolver, error) {
    wxUser, err := rp.L("user").(*rp.UserRepository).WxLogin(ctx, args.Code)
    if err != nil {
        return nil, err
    }
    user, err := rp.L("user").(*rp.UserRepository).FindByID(wxUser.UserID)
    if err != nil {
        return nil, err
    }
    user.Openid = wxUser.OpenId
    rp.L("user").(*rp.UserRepository).UpdateLoginStatus(ctx, user.ID)
    return &userResolver{user}, nil
}

func (r *Resolver) MyCoupons(ctx context.Context, args struct{ Status string }) ([]*couponResolver, error) {
    c, err := rp.L("coupon").(*rp.CouponRepository).FindByStatus(ctx, args.Status, 100, nil)
    if err != nil {
        return nil, err
    }
    res := make([]*couponResolver, len(c))
    for i, v := range c {
        res[i] = &couponResolver{v}
    }
    return res, nil
}

func (r *Resolver) Coupons(ctx context.Context, args struct {
    First  *int32
    After  *string
    Status string
}) (*couponsConnectionResolver, error) {
    fetchSize := int(util.GetInt32(args.First, DefaultPageSize))
    decodedIndex, _ := service.DecodeCursor(args.After)
    list, err := rp.L("coupon").(*rp.CouponRepository).FindByStatus(ctx, args.Status, fetchSize, decodedIndex)
    if err != nil {
        return nil, err
    }
    count, err := rp.L("coupon").(*rp.CouponRepository).Count(ctx, args.Status)
    if err != nil {
        return nil, err
    }
    var from *string
    var to *string
    if len(list) > 0 {
        from = &(list[0].ID)
        to = &(list[len(list)-1].ID)
    }
    return &couponsConnectionResolver{list, Connection{totalCount: count, from: from, to: to, hasNext: len(list) == int(fetchSize)}}, nil
}
