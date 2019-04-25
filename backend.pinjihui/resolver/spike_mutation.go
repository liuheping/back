package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CreateSpike(ctx context.Context, args struct {
	Spike *model.Spike
}) (*spikeResolver, error) {
	sp, err := rp.L("spike").(*rp.SpikeRepository).SaveSpike(ctx, args.Spike)
	if err != nil {
		return nil, err
	}
	return &spikeResolver{sp}, nil
}

func (r *Resolver) UpdateSpike(ctx context.Context, args struct {
	ID    graphql.ID
	Spike *model.Spike
}) (*spikeResolver, error) {
	args.Spike.ID = string(args.ID)
	sp, err := rp.L("spike").(*rp.SpikeRepository).SaveSpike(ctx, args.Spike)
	if err != nil {
		return nil, err
	}
	return &spikeResolver{sp}, nil
}
