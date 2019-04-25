package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) FindAdPosition(ctx context.Context) (*string, error) {
	options, err := rp.L("ad").(*rp.AdRepository).FindAdPositionOptions(ctx)
	if err != nil {
		return nil, err
	}
	return options, nil
}

func (r *Resolver) Ad(ctx context.Context, args struct {
	ID string
}) (*adResolver, error) {
	ad, err := rp.L("ad").(*rp.AdRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &adResolver{ad}, nil
}

// 获取所有广告
func (r *Resolver) Ads(ctx context.Context) (*[]*adResolver, error) {
	ads, err := rp.L("ad").(*rp.AdRepository).FindAll(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*adResolver, len(ads))
	for i := range l {
		l[i] = &adResolver{(ads)[i]}
	}
	return &l, nil
}
