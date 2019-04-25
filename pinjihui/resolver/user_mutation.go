package resolver

import (
    "pinjihui.com/pinjihui/model"
    rp "pinjihui.com/pinjihui/repository"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
)

func (r *Resolver) Register(ctx context.Context, args *struct {
    Mobile   string
    Password string
}) (*userResolver, error) {
    defaultName := model.DefaultMobileUserNickName
    user := &model.User{
        Mobile:   &args.Mobile,
        Password: args.Password,
        LastIp:   *ctx.Value("requester_ip").(*string),
        Type:     model.UTConsumer,
        Name:     &defaultName,
    }

    user, err := rp.L("user").(*rp.UserRepository).CreateUser(user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }
    ctx.Value("log").(*logging.Logger).Debugf("Created user : %v", *user)
    return &userResolver{user}, nil
}

func (r *Resolver) SaveWxUserInfo(ctx context.Context, args *struct {
    WxUserInfo *rp.WxUser
}) (bool, error) {
    return rp.L("user").(*rp.UserRepository).SaveWxUserInfo(ctx, args.WxUserInfo)
}

func (r *Resolver) ReceiveCoupon(ctx context.Context, args struct{ InviteCode string }) (bool, error) {
    success, err := rp.L("user").(*rp.UserRepository).ReceiveCoupon(ctx, args.InviteCode)
    if err != nil {
        return false, err
    }
    return success, nil
}

func (r *Resolver) ReceiveShareCoupon(ctx context.Context, args struct{ Type string }) (success bool, err error) {
    err = rp.L("user").(*rp.UserRepository).AddShareCoupon(ctx, args.Type)
    if err == nil {
        success = true
    }
    return
}

func (r *Resolver) BindPhoneNumber(ctx context.Context, args struct{ Number string }) (bool, error) {
    err := rp.L("user").(*rp.UserRepository).BindPhoneNumber(ctx, args.Number)
    if err != nil {
        return false, err
    }
    return true, nil
}
