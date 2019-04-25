package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
)

type paymentResolver struct {
	m *model.Payment
}

func (r *paymentResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *paymentResolver) PayName() string {
	return r.m.Pay_name
}

func (r *paymentResolver) PayCode() string {
	return r.m.Pay_code
}

func (r *paymentResolver) PayFee() *string {
	return r.m.Pay_fee
}

func (r *paymentResolver) PayDesc() *string {
	return r.m.Pay_desc
}

func (r *paymentResolver) SortOrder() *int32 {
	return r.m.Sort_order
}

func (r *paymentResolver) Enabled() bool {
	return r.m.Enabled
}

func (r *paymentResolver) IsOnline() bool {
	return r.m.Is_online
}

func (r *paymentResolver) Iscod() bool {
	return r.m.Is_cod
}
