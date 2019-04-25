package resolver

import (
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/loader"
    rp "pinjihui.com/pinjihui/repository"
    "golang.org/x/net/context"
)

type cartItemGroupResolver struct {
    items      []*model.CartItem
    merchantID string
}

func (c *cartItemGroupResolver) Items() ([]*cartItemResolver, error) {
    itemResolvers := make([]*cartItemResolver, len(c.items))
    for i, v := range c.items {
        itemResolvers[i] = &cartItemResolver{v}
    }
    return itemResolvers, nil
}

func (c *cartItemGroupResolver) Merchant(ctx context.Context) (*merchantResolver, error) {
    merchant, err := loader.LoadMerchant(ctx, c.merchantID)

    if err != nil {
        return nil, err
    }
    return &merchantResolver{merchant}, nil
}

func (c *cartItemGroupResolver) ShippingMethods() ([]*shippingMethodResolver, error) {
    sms, err := rp.L("shipping").(*rp.ShippingMethodRepository).FindAll()
    if err != nil {
        return nil, err
    }
    res := make([]*shippingMethodResolver, 0, len(sms))
    for _, v := range sms {
        if v.EnableForPlatform == false && c.merchantID == rp.PLATFORM {
            continue
        }
        res = append(res, &shippingMethodResolver{v})
    }
    return res, nil
}
