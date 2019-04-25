package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
)

type takeCashResolver struct {
	u *model.TakeCash
}

func (r *takeCashResolver) UserID() *graphql.ID {
	a := graphql.ID(r.u.UserID)
	return &a
}

func (r *takeCashResolver) IsChecked() *bool {
	return &r.u.IsChecked
}

func (r *takeCashResolver) DebitCardInfo() *debitCardInfoResolver {
	res := debitCardInfoResolver{r.u.DebitCardInfo}
	return &res
}
