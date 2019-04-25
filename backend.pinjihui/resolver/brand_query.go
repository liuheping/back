package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) Brand(ctx context.Context, args struct {
	ID string
}) (*brandResolver, error) {
	brand, err := rp.L("brand").(*rp.BrandRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &brandResolver{brand}, nil
}

//获取所有品牌
func (r *Resolver) Brands(ctx context.Context) (*[]*brandResolver, error) {
	brand, err := rp.L("brand").(*rp.BrandRepository).FindAll(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*brandResolver, len(brand))
	for i := range l {
		l[i] = &brandResolver{(brand)[i]}
	}
	return &l, nil
}

// 根据机型查找品牌
func (r *Resolver) BrandByMachine(ctx context.Context, args struct {
	MachineType string
}) (*brandResolver, error) {
	brand, err := rp.L("brand").(*rp.BrandRepository).FindByMachineType(ctx, args.MachineType)
	if err != nil {
		return nil, err
	}
	return &brandResolver{brand}, nil
}
