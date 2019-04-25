package resolver

import (
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) Config(ctx context.Context, args struct {
	ID string
}) (*configResolver, error) {
	con, err := rp.L("config").(*rp.ConfigRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &configResolver{con}, nil
}

// 获取所有配置
func (r *Resolver) Configs(ctx context.Context) (*[]*configResolver, error) {
	cons, err := rp.L("config").(*rp.ConfigRepository).FindAll(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*configResolver, len(cons))
	for i := range l {
		l[i] = &configResolver{(cons)[i]}
	}
	return &l, nil
}
