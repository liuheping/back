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

func (r *Resolver) CreateAttribute(ctx context.Context, args struct {
	NewAttr *model.AttributeARR
}) (*attributeResolver, error) {
	gc.CheckAuth(ctx)
	if args.NewAttr.Options != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewAttr.Options, "\",\""))
		args.NewAttr.Attribute.Options = &str
	}
	args.NewAttr.Attribute.Merchant_id = gc.CurrentUser(ctx)

	attribute, err := rp.L("attribute").(*rp.AttributeRepository).SaveAttribute(ctx, &args.NewAttr.Attribute)
	if err != nil {
		return nil, err
	}
	return &attributeResolver{attribute}, nil
}

func (r *Resolver) UpdateAttribute(ctx context.Context, args struct {
	ID      graphql.ID
	NewAttr *model.AttributeARR
}) (*attributeResolver, error) {
	gc.CheckAuth(ctx)
	if args.NewAttr.Options != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewAttr.Options, "\",\""))
		args.NewAttr.Attribute.Options = &str
	}
	args.NewAttr.Attribute.ID = string(args.ID)
	attribute, err := rp.L("attribute").(*rp.AttributeRepository).SaveAttribute(ctx, &args.NewAttr.Attribute)
	if err != nil {
		return nil, err
	}
	return &attributeResolver{attribute}, nil
}

func (r *Resolver) DeleteAttribute(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("attribute").(*rp.AttributeRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
