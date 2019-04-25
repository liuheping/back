package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

//根据ID查找属性集
func (r *Resolver) AttributeSet(ctx context.Context, args struct {
	ID string
}) (*attributeSetResolver, error) {
	attr, err := rp.L("attributeset").(*rp.AttributeSetRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &attributeSetResolver{attr}, nil
}

//获取所有属性集合
func (r *Resolver) AttributeSets(ctx context.Context) (*[]*attributeSetResolver, error) {
	sets, err := rp.L("attributeset").(*rp.AttributeSetRepository).FindAll(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*attributeSetResolver, len(sets))
	for i := range l {
		l[i] = &attributeSetResolver{(sets)[i]}
	}
	return &l, nil
}
