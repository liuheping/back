package resolver

import (
    "golang.org/x/net/context"
    "github.com/graph-gophers/graphql-go"
    rp "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/util"
    "pinjihui.com/pinjihui/service"
    "strconv"
)

func (r *Resolver) Checkout(ctx context.Context, args struct{ Ids []graphql.ID }) (*checkoutResolver, error) {
    items, err := rp.L("cart").(*rp.CartRepository).FindByIDs(ctx, ID2String(&args.Ids))
    if err != nil {
        return nil, err
    }
    return &checkoutResolver{items}, nil
}

func ID2String(ids *[]graphql.ID) (*[]string) {
    if ids == nil {
        return nil
    }
    idStrs := make([]string, len(*ids))
    for i, v := range *ids {
        idStrs[i] = string(v)
    }
    return &idStrs
}

func (r *Resolver) Orders(ctx context.Context, args struct {
    Status *string
    First  *int32
    After  *string
}) (*ordersConnectionResolver, error) {
    fetchSize := int(util.GetInt32(args.First, DefaultPageSize))
    decodedIndex, _ := service.DecodeCursor(args.After)
    list, err := rp.L("order").(*rp.OrderRepository).List(ctx, args.Status, fetchSize, decodedIndex)
    if err != nil {
        return nil, err
    }
    count, err := rp.L("order").(*rp.OrderRepository).Count(ctx, args.Status)
    if err != nil {
        return nil, err
    }
    var from *string
    var to *string
    if len(list) > 0 {
        from = &(list[0].ID)
        to = &(list[len(list)-1].ID)
    }
    return &ordersConnectionResolver{list, Connection{totalCount: count, from: from, to: to, hasNext: len(list) == int(fetchSize)}}, nil
}

func (r *Resolver) Order(ctx context.Context, args struct{ ID graphql.ID }) (*orderResolver, error) {
    order, err := rp.L("order").(*rp.OrderRepository).FindByID(ctx, string(args.ID))
    if err != nil {
        return nil, err
    }
    return &orderResolver{m: order}, nil
}

func (r *Resolver) MyProducts(ctx context.Context, args struct {
    First        *int32
    After        *string
    HasCommented bool
}) (*OrderProductItemConnnectionResolver, error) {
    fetchSize, offset := getPageParams(args.First, args.After)
    list, err := rp.L("order").(*rp.OrderRepository).MyOrderProducts(ctx, fetchSize, offset, args.HasCommented)
    if err != nil {
        return nil, err
    }
    count, err := rp.L("order").(*rp.OrderRepository).MyOrderProductsCount(ctx, args.HasCommented)
    if err != nil {
        return nil, err
    }
    var from *string
    var to *string

    nOffset := strconv.Itoa(int(offset) + len(list))
    to = &nOffset
    return &OrderProductItemConnnectionResolver{list, Connection{totalCount: count, from: from, to: to, hasNext: len(list) == int(fetchSize)}}, nil
}
