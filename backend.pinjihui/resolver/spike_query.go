package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) Spike(ctx context.Context, args struct {
	ID string
}) (*spikeResolver, error) {
	sp, err := rp.L("spike").(*rp.SpikeRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &spikeResolver{sp}, nil
}

// 获取所有配置
func (r *Resolver) Spikes(ctx context.Context) (*[]*spikeResolver, error) {
	sps, err := rp.L("spike").(*rp.SpikeRepository).FindAll(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*spikeResolver, len(sps))
	for i := range l {
		l[i] = &spikeResolver{(sps)[i]}
	}
	return &l, nil
}
