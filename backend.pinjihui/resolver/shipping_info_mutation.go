package resolver

import (
	"fmt"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CreateShippingInfo(ctx context.Context, args struct {
	Info *model.ShippingInfoARR
}) (*shippingInfoResolver, error) {
	gc.CheckAuth(ctx)
	if args.Info.Images != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.Info.Images, "\",\""))
		args.Info.ShippingInfo.Images = &str
	}
	shinningInfo, err := rp.L("shippinginfo").(*rp.ShippingInfoRepository).SaveShippingInfo(ctx, &args.Info.ShippingInfo)
	if err != nil {
		return nil, err
	}
	return &shippingInfoResolver{shinningInfo}, nil
}

func (r *Resolver) UpdateShippingInfo(ctx context.Context, args struct {
	ID   graphql.ID
	Info *model.ShippingInfoARR
}) (*shippingInfoResolver, error) {
	gc.CheckAuth(ctx)
	if args.Info.Images != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.Info.Images, "\",\""))
		args.Info.ShippingInfo.Images = &str
	}
	args.Info.ShippingInfo.ID = string(args.ID)
	shinningInfo, err := rp.L("shippinginfo").(*rp.ShippingInfoRepository).SaveShippingInfo(ctx, &args.Info.ShippingInfo)
	if err != nil {
		return nil, err
	}
	return &shippingInfoResolver{shinningInfo}, nil
}

func (r *Resolver) DeleteShippingInfo(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("shippinginfo").(*rp.ShippingInfoRepository).DeletedByID(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
