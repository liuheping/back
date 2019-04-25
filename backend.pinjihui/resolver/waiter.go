package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
)

type waiterResolver struct {
	m *model.Waiter
}

func (r *waiterResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *waiterResolver) Merchant(ctx context.Context) (*merchantProfileResolver, error) {
	if r.m.Merchant_id == nil {
		return nil, nil
	}
	mp, err := ctx.Value("userRepository").(*repository.UserRepository).MerchantProfile(*r.m.Merchant_id)
	if err != nil {
		return nil, err
	}
	return &merchantProfileResolver{mp}, nil
}

func (r *waiterResolver) Mobile() string {
	return r.m.Mobile
}

func (r *waiterResolver) Name() *string {
	return r.m.Name
}

func (r *waiterResolver) Waiter_id() *string {
	return r.m.Waiter_id
}

func (r *waiterResolver) Checked() bool {
	return r.m.Checked
}

func (r *waiterResolver) Deleted() bool {
	return r.m.Deleted
}

func (r *waiterResolver) Remark() *string {
	return r.m.Remark
}
