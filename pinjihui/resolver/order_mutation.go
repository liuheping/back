package resolver

import (
    "golang.org/x/net/context"
    rp "pinjihui.com/pinjihui/repository"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
)

type OrderInput struct {
    ShippingAddrID graphql.ID
    ItemID         []graphql.ID
    SmOpt []*struct {
        MerchantID     graphql.ID
        ShippingMethod string
        Message        *string
    }
    UsedCoupon *[]string
    InvOpt     *model.InvOpt
}

func (r *Resolver) CreateOrder(ctx context.Context, args *OrderInput) (*orderResolver, error) {
    order, err := rp.L("order").(*rp.OrderRepository).Create(ctx, args2OrderInput(args))
    if err != nil {
        return nil, err
    }
    return &orderResolver{m: order[0]}, nil
}

func args2OrderInput(input *OrderInput) *model.OrderInput {
    smOpt := make([]*model.ShppingMethodOption, len(input.SmOpt))
    for i, v := range input.SmOpt {
        smOpt[i] = &model.ShppingMethodOption{string(v.MerchantID), string(v.ShippingMethod), v.Message}
    }
    return &model.OrderInput{
        ShippingAddrID: string(input.ShippingAddrID),
        ItemIDs:        *ID2String(&input.ItemID),
        SmOpt:          smOpt,
        UsedCoupon:     input.UsedCoupon,
        InvOpt:         input.InvOpt,
    }
}

func (r *Resolver) CancelOrder(ctx context.Context, args struct{ ID string }) (bool, error) {
    err := rp.L("order").(*rp.OrderRepository).Cancel(ctx, args.ID)
    if err != nil {
        return false, err
    }
    return true, nil
}
func (r *Resolver) ConfirmReceipt(ctx context.Context, args struct{ ID string }) (bool, error) {
    err := rp.L("order").(*rp.OrderRepository).ConfirmReceipt(ctx, args.ID)
    if err != nil {
        return false, err
    }
    return true, nil
}
