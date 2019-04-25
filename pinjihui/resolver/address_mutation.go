package resolver

import (
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/model"
    rp "pinjihui.com/pinjihui/repository"
    "github.com/graph-gophers/graphql-go"
)

func (r *Resolver) CreateAddress(ctx context.Context, args struct {
    NewAddr *model.ShippingAddress
}) (*shippingAddressResolver, error) {
    addr, err := rp.L("address").(*rp.AddressRepository).Save(ctx, args.NewAddr)
    if err != nil {
        return nil, err
    }
    return &shippingAddressResolver{addr}, nil
}

func (r *Resolver) UpdateAddress(ctx context.Context, args struct {
    ID      graphql.ID
    NewAddr *model.ShippingAddress
}) (*shippingAddressResolver, error) {
    args.NewAddr.ID = string(args.ID)
    addr, err := rp.L("address").(*rp.AddressRepository).Save(ctx, args.NewAddr)
    if err != nil {
        return nil, err
    }
    return &shippingAddressResolver{addr}, nil
}

func (r *Resolver) SetDefaultAddress(ctx context.Context, args struct{ ID graphql.ID }) (bool, error) {
    return rp.L("address").(*rp.AddressRepository).SetDefault(ctx, string(args.ID))
}

func (r *Resolver) DeleteAddress(ctx context.Context, args struct{ ID graphql.ID }) (bool, error) {
    return rp.L("address").(*rp.AddressRepository).Delete(ctx, string(args.ID))
}
