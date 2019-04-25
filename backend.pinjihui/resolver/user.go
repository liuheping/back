package resolver

import (
	"time"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/service"
	//"pinjihui.com/backend.pinjihui/service"
)

type userResolver struct {
	u *model.User
}

func (r *userResolver) ID() graphql.ID {
	return graphql.ID(r.u.ID)
}

func (r *userResolver) Name() *string {
	return r.u.Name
}

func (r *userResolver) Mobile() *string {
	return r.u.Mobile
}

func (r *userResolver) Email() *string {
	return r.u.Email
}

func (r *userResolver) LastIp() *string {
	return r.u.LastIp
}

func (r *userResolver) CreatedAt() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.u.CreatedAt)
	return &graphql.Time{Time: res}, err
}

func (r *userResolver) UpdatedAt() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.u.UpdatedAt)
	return &graphql.Time{Time: res}, err
}

func (r *userResolver) LastLoginTime() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, *r.u.LastLoginTime)
	return &graphql.Time{Time: res}, err
}

func (r *userResolver) Type() *string {
	return &r.u.Type
}

func (r *userResolver) Status() string {
	return r.u.Status
}

func (r *userResolver) UserProfile(ctx context.Context) (*userProfileResolver, error) {
	var a interface{}
	if r.u.Type != "consumer" {
		mp, err := ctx.Value("userRepository").(*rp.UserRepository).MerchantProfile(r.u.ID)
		if err != nil {
			return nil, err
		}
		if mp == nil {
			return nil, nil
		}
		a = &merchantProfileResolver{mp}
	} else {
		cp, err := rp.L("wechartprofile").(*rp.WechartProfileRepository).FindByUserID(ctx, r.u.ID)
		if err != nil {
			return nil, err
		}
		if cp == nil {
			return nil, nil
		}
		a = &customerProfileResolver{cp}
	}
	return &userProfileResolver{a}, nil
}

func (r *userResolver) LastOrders() *[]*orderResolver {
	res := make([]*orderResolver, 3)
	for i := range res {
		v := orderResolver{}
		res[i] = &v
	}
	return &res
}

func (r *userResolver) Tokenstring(ctx context.Context) (*string, error) {
	return ctx.Value("authService").(*service.AuthService).SignJWT(r.u)
}
