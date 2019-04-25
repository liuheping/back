package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CreateConfig(ctx context.Context, args struct {
	NewCon *model.Config
}) (*configResolver, error) {
	con, err := rp.L("config").(*rp.ConfigRepository).SaveConfig(ctx, args.NewCon)
	if err != nil {
		return nil, err
	}
	return &configResolver{con}, nil
}

func (r *Resolver) UpdateConfig(ctx context.Context, args struct {
	ID     graphql.ID
	NewCon *model.Config
}) (*configResolver, error) {
	args.NewCon.ID = string(args.ID)
	con, err := rp.L("config").(*rp.ConfigRepository).SaveConfig(ctx, args.NewCon)
	if err != nil {
		return nil, err
	}
	return &configResolver{con}, nil
}

func (r *Resolver) DeleteConfig(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("config").(*rp.ConfigRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
