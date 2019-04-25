package resolver

import (
    "pinjihui.com/pinjihui/model"
    rp "pinjihui.com/pinjihui/repository"
    "golang.org/x/net/context"
	"pinjihui.com/pinjihui/util"
)

type checkoutResolver struct {
    items []*model.CartItem
}

func (c *checkoutResolver) Groups() ([]*cartItemGroupResolver, error) {
    merchantMapItems := make(map[string][]*model.CartItem)
    for _, v := range c.items {
        merchantMapItems[v.MerchantID] = append(merchantMapItems[v.MerchantID], v)
    }
    itemResolvers := make([]*cartItemGroupResolver, len(merchantMapItems))
    i := 0
    for k, v := range merchantMapItems {
        itemResolvers[i] = &cartItemGroupResolver{v, k}
        i++
    }
    return itemResolvers, nil
}

func (c *checkoutResolver) ProductAmount() string {
    amount := 0.0
    for _, item := range c.items {
        amount += float64(item.ProductCount) * item.Price
    }
    return util.FmtMoney(amount)
}

func (c *checkoutResolver) Coupons(ctx context.Context) ([]*couponResolver, error) {
    cs, err := rp.L("coupon").(*rp.CouponRepository).FindAvailable(ctx, 100, nil)
    if err != nil {
        return nil, err
    }

    res := make([]*couponResolver, len(cs))
    for i, v := range cs {
       res[i] = &couponResolver{v}
    }
    return res, nil
}
