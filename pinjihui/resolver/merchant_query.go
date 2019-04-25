package resolver

import (
    "golang.org/x/net/context"
    rp "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/loader"
    "pinjihui.com/pinjihui/model"
    "strconv"
    "github.com/graph-gophers/graphql-go"
)

func (r *Resolver) Merchant(ctx context.Context, args struct{ Id string }) (*merchantResolver, error) {
    merchant, err := loader.LoadMerchant(ctx, args.Id)
    if err != nil {
        return nil, err
    }
    return &merchantResolver{merchant}, nil
}

func (r *Resolver) Merchants(ctx context.Context, args struct {
    First    *int32
    After    *string
    OrderBy  string
    Position *model.Location
}) (*merchantsConnectionResolver, error) {
    fetchSize, offset := getPageParams(args.First, args.After)
    list, err := rp.L("merchant").(*rp.MerchantRepository).List(fetchSize, offset, args.OrderBy, args.Position)
    if err != nil {
        return nil, err
    }
    count, err := rp.L("merchant").(*rp.MerchantRepository).Count()
    if err != nil {
        return nil, err
    }
    var from *string
    var to *string

    nOffset := strconv.Itoa(int(offset) + len(list))
    to = &nOffset
    return &merchantsConnectionResolver{m: list, totalCount: count, from: from, to: to, hasNext: len(list) == int(fetchSize)}, nil
}

func (r *Resolver) AgentPurchase(ctx context.Context, args struct{ ProductID graphql.ID }) (bool, error) {
    err := rp.L("merchant").(*rp.MerchantRepository).Purchase(ctx, string(args.ProductID))
    if err != nil {
        return false, err
    }
    return true, nil
}
func (r *Resolver) AgentUnPurchase(ctx context.Context, args struct{ ProductID graphql.ID }) (bool, error) {
    err := rp.L("merchant").(*rp.MerchantRepository).UnPurchase(ctx, string(args.ProductID))
    if err != nil {
        return false, err
    }
    return true, nil
}
