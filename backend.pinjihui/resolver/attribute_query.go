package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) Attribute(ctx context.Context, args struct {
	ID string
}) (*attributeResolver, error) {
	attribute, err := rp.L("attribute").(*rp.AttributeRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &attributeResolver{attribute}, nil
}

//获取所有属性
func (r *Resolver) Attributes(ctx context.Context) (*[]*attributeResolver, error) {
	attr, err := rp.L("attribute").(*rp.AttributeRepository).FindAll(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*attributeResolver, len(attr))
	for i := range l {
		l[i] = &attributeResolver{(attr)[i]}
	}
	return &l, nil
}

//根据ID集合获取所有属性
func (r *Resolver) AttributesByIds(ctx context.Context, args struct{ Ids *[]string }) (*[]*attributeResolver, error) {
	attrs, err := rp.L("attribute").(*rp.AttributeRepository).FindByIDs(ctx, args.Ids)
	if err != nil {
		return nil, err
	}
	l := make([]*attributeResolver, len(attrs))
	for i := range l {
		l[i] = &attributeResolver{(attrs)[i]}
	}
	return &l, nil
}
