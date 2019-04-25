package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

// 分类列表, 顶级
func (r *Resolver) Categories() (*[]*categoryResolver, error) {
	return GetChildResolver(nil)
}

func (r *Resolver) Category(ctx context.Context, args struct {
	ID string
}) (*categoryResolver, error) {
	cat, err := rp.L("category").(*rp.CategoryRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &categoryResolver{cat}, nil
}
