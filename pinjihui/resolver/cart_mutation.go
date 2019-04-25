package resolver

import (
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/model"
    rp "pinjihui.com/pinjihui/repository"
    "github.com/graph-gophers/graphql-go"
)

func (r *Resolver) AddToCart(ctx context.Context, args struct {
    Item *model.CartItem
}) (*cartItemResolver, error) {
    newItem, err := rp.L("cart").(*rp.CartRepository).Save(ctx, args.Item)
    if err != nil {
        return nil, err
    }
    return &cartItemResolver{newItem}, nil
}

func (r *Resolver) UpdateCountInCart(ctx context.Context, args struct {
    ID    graphql.ID
    Count int32
}) (*cartItemResolver, error) {
    newItem, err := rp.L("cart").(*rp.CartRepository).UpdateCount(ctx, string(args.ID), args.Count)
    if err != nil {
        return nil, err
    }
    return &cartItemResolver{newItem}, nil
}

func (r *Resolver) DeleteCartItem(ctx context.Context, args struct{ID graphql.ID}) (bool, error) {
    return rp.L("cart").(*rp.CartRepository).DeleteItem(ctx, string(args.ID))
}
