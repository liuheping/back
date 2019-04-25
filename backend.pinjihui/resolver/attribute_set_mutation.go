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

func (r *Resolver) CreateAttributeSet(ctx context.Context, args struct {
	NewAttrSet *model.AttributeSetARR
}) (*attributeSetResolver, error) {
	gc.CheckAuth(ctx)

	if args.NewAttrSet.Attribute_ids != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewAttrSet.Attribute_ids, "\",\""))
		args.NewAttrSet.AttributeSet.Attribute_ids = &str
	}
	args.NewAttrSet.AttributeSet.Merchant_id = gc.CurrentUser(ctx)
	attributeset, err := rp.L("attributeset").(*rp.AttributeSetRepository).SaveAttributeSet(ctx, &args.NewAttrSet.AttributeSet)
	if err != nil {
		return nil, err
	}
	return &attributeSetResolver{attributeset}, nil
}

func (r *Resolver) UpdateAttributeSet(ctx context.Context, args struct {
	ID         graphql.ID
	NewAttrSet *model.AttributeSetARR
}) (*attributeSetResolver, error) {
	gc.CheckAuth(ctx)
	if args.NewAttrSet.Attribute_ids != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewAttrSet.Attribute_ids, "\",\""))
		args.NewAttrSet.AttributeSet.Attribute_ids = &str
	}
	args.NewAttrSet.AttributeSet.ID = string(args.ID)
	attributeset, err := rp.L("attributeset").(*rp.AttributeSetRepository).SaveAttributeSet(ctx, &args.NewAttrSet.AttributeSet)
	if err != nil {
		return nil, err
	}
	return &attributeSetResolver{attributeset}, nil
}

func (r *Resolver) DeleteAttributeSet(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("attributeset").(*rp.AttributeSetRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
