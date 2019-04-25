package resolver

import (
	"time"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
)

type cashRequestResolver struct {
	u *model.CashRequest
}

func (r *cashRequestResolver) ID() graphql.ID {
	return graphql.ID(r.u.ID)
}

// func (r *cashRequestResolver) MerchantID() graphql.ID {
// 	return graphql.ID(r.u.Merchant_id)
// }

func (r *cashRequestResolver) Merchant(ctx context.Context) (*merchantProfileResolver, error) {
	mp, err := ctx.Value("userRepository").(*repository.UserRepository).MerchantProfile(r.u.Merchant_id)
	if err != nil {
		return nil, err
	}
	return &merchantProfileResolver{mp}, nil
}

func (r *cashRequestResolver) Amount() float64 {
	return r.u.Amount
}

func (r *cashRequestResolver) DebitCardInfo() *debitCardInfoResolver {
	res := debitCardInfoResolver{r.u.DebitCardInfo}
	return &res
}

func (r *cashRequestResolver) Status() string {
	return r.u.Status
}

func (r *cashRequestResolver) Reply() *string {
	return r.u.Reply
}

func (r *cashRequestResolver) Note() *string {
	return r.u.Note
}

func (r *cashRequestResolver) CreatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.u.Created_at)
	return graphql.Time{Time: res}, err
}

func (r *cashRequestResolver) UpdatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.u.Updated_at)
	return graphql.Time{Time: res}, err
}
