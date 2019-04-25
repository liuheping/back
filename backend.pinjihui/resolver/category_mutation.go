package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CreateCategory(ctx context.Context, args struct {
	NewCat *model.Category
}) (*categoryResolver, error) {
	cat, err := rp.L("category").(*rp.CategoryRepository).SaveCategory(ctx, args.NewCat)
	if err != nil {
		return nil, err
	}
	return &categoryResolver{cat}, nil
}

func (r *Resolver) UpdateCategory(ctx context.Context, args struct {
	ID     graphql.ID
	NewCat *model.Category
}) (*categoryResolver, error) {
	args.NewCat.ID = string(args.ID)
	cat, err := rp.L("category").(*rp.CategoryRepository).SaveCategory(ctx, args.NewCat)
	if err != nil {
		return nil, err
	}
	return &categoryResolver{cat}, nil
}

func (r *Resolver) DeleteCategory(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("category").(*rp.CategoryRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
