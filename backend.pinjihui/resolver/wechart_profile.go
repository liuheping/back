package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/loader"
	"pinjihui.com/backend.pinjihui/model"
)

type wechartProfileResolver struct {
	m *model.WechartProfile
}

func (r *wechartProfileResolver) Openid() graphql.ID {
	return graphql.ID(r.m.Openid)
}

func (r *wechartProfileResolver) Nickname() string {
	return r.m.Nick_name
}

func (r *wechartProfileResolver) Gender() int32 {
	return r.m.Gender
}

func (r *wechartProfileResolver) Language() *string {
	return r.m.Language
}

func (r *wechartProfileResolver) City() *string {
	return r.m.City
}

func (r *wechartProfileResolver) Province() *string {
	return r.m.Province
}

func (r *wechartProfileResolver) County() *string {
	return r.m.Country
}

func (r *wechartProfileResolver) AvatarUrl() *string {
	return r.m.Avatar_url
}

func (r *wechartProfileResolver) User(ctx context.Context) (*userResolver, error) {
	user, err := loader.LoadUser(ctx, r.m.User_id)
	if err != nil {
		return nil, err
	}
	return &userResolver{user}, nil
}
