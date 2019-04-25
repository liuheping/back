package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CreateAd(ctx context.Context, args struct {
	Ad *model.Ad
}) (*adResolver, error) {
	ad, err := rp.L("ad").(*rp.AdRepository).SaveAd(ctx, args.Ad)
	if err != nil {
		return nil, err
	}
	return &adResolver{ad}, nil
}

func (r *Resolver) UpdateAd(ctx context.Context, args struct {
	ID graphql.ID
	Ad *model.Ad
}) (*adResolver, error) {
	args.Ad.ID = string(args.ID)
	ad, err := rp.L("ad").(*rp.AdRepository).SaveAd(ctx, args.Ad)
	if err != nil {
		return nil, err
	}
	return &adResolver{ad}, nil
}

func (r *Resolver) DeleteAd(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("ad").(*rp.AdRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Resolver) SetAdIsShow(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("ad").(*rp.AdRepository).Isshow(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
