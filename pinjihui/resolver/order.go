package resolver

import (
    "pinjihui.com/pinjihui/model"
    "time"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/loader"
    "golang.org/x/net/context"
    rp "pinjihui.com/pinjihui/repository"
    gc "pinjihui.com/pinjihui/context"
    "pinjihui.com/pinjihui/util"
)

type orderResolver struct {
    m *model.Order
}

func (r *orderResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *orderResolver) OrderStatus() string {
    return r.m.Status
}

func (r *orderResolver) Children() (*[]*orderResolver, error) {
    orders, err := r.m.GetChildren(rp.L("order").(*rp.OrderRepository).FindChildren)
    if err != nil {
        return nil, err
    }
    res := make([]*orderResolver, len(orders))
    for i, v := range orders {
        res[i] = &orderResolver{m: v}
    }
    return &res, nil
}

func (r *orderResolver) Address() (*orderAddressResolver, error) {
    if err := r.m.ParseAddress(); err != nil {
        return nil, err
    }
    return &orderAddressResolver{r.m.Address}, nil
}

func (r *orderResolver) Merchant(ctx context.Context) (*merchantResolver, error) {
    if r.m.MerchantID == nil {
        return nil, nil
    }
    merchant, err := loader.LoadMerchant(ctx, *r.m.MerchantID)
    if err != nil {
        return nil, err
    }
    return &merchantResolver{merchant}, nil
}

func (r *orderResolver) Products() (*[]*orderProductItemResolver, error) {
    items, err := r.m.GetProducts(rp.L("order").(*rp.OrderRepository).FindOrderProducts)
    if err != nil {
        return nil, err
    }
    res := make([]*orderProductItemResolver, len(items))
    for i, v := range items {
        res[i] = &orderProductItemResolver{v}
    }
    return &res, nil
}

func (r *orderResolver) ShippingName() *string {
    return r.m.ShippingName
}

func (r *orderResolver) PayName() string {
    return r.m.PayName
}

func (r *orderResolver) Amount() string {
    return util.FmtMoney(r.m.Amount)
}

func (r *orderResolver) OrderAmount() string {
    return util.FmtMoney(r.m.OrderAmount)
}

func (r *orderResolver) OfferAmount() string {
    return util.FmtMoney(r.m.OfferAmount)
}

func (r *orderResolver) CreatedAt() (graphql.Time, error) {
    res, err := time.Parse(time.RFC3339, r.m.CreatedAt)
    return graphql.Time{Time: res}, err
}

func (r *orderResolver) PayTime() (*graphql.Time, error) {
    if r.m.PayTime == nil {
        return nil, nil
    }
    res, err := time.Parse(time.RFC3339, *r.m.PayTime)
    return &graphql.Time{Time: res}, err
}

func (r *orderResolver) ShippingTime() (*graphql.Time, error) {
    if r.m.ShippingTime == nil {
        return nil, nil
    }
    res, err := time.Parse(time.RFC3339, *r.m.ShippingTime)
    return &graphql.Time{Time: res}, err
}

func (r *orderResolver) ShippingInfo(ctx context.Context) (*shippingInfoResolver, error) {
    info, err := rp.L("order").(*rp.OrderRepository).FindShippingInfo(ctx, r.m.ID)
    if err == gc.ErrNoRecord {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &shippingInfoResolver{info}, nil
}

func (r *orderResolver) WechartPayParams(ctx context.Context) (*wechartPayParamsResolver, error) {
    if r.m.Status != model.OS_UNPAID || r.m.ParentID != nil {
        return nil, nil
    }
    p, err := rp.L("order").(*rp.OrderRepository).WechartPayPrepare(ctx, r.m)
    if err != nil {
        return nil, err
    }
    return &wechartPayParamsResolver{p}, nil
}

func (r *orderResolver) ProductQuantity() (int32, error) {
    var count int32 = 0
    items, err := r.m.GetProducts(rp.L("order").(*rp.OrderRepository).FindOrderProducts)
    if err != nil {
        return 0, err
    }
    if len(items) == 0 && r.m.ParentID == nil {
        //取子订单
        orders, err := r.m.GetChildren(rp.L("order").(*rp.OrderRepository).FindChildren)
        if err != nil {
            return 0, err
        }
        for _, o := range orders {
            items, err = o.GetProducts(rp.L("order").(*rp.OrderRepository).FindOrderProducts)
            if err != nil {
                return 0, err
            }
            for _, p := range items {
                count += p.ProductNumber
            }
        }
    } else {
        for _, p := range items {
            count += p.ProductNumber
        }
    }
    return count, nil
}
