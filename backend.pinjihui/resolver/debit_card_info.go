package resolver

import "pinjihui.com/backend.pinjihui/model"

type debitCardInfoResolver struct {
	m *model.DebitCardInfo
}

func (r *debitCardInfoResolver) CardHolder() string {
	return r.m.CardHolder
}

func (r *debitCardInfoResolver) BankName() string {
	return r.m.BankName
}

func (r *debitCardInfoResolver) CardNumber() string {
	return r.m.CardNumber
}

func (r *debitCardInfoResolver) Province() int32 {
	return r.m.ProvinceID
}

func (r *debitCardInfoResolver) City() int32 {
	return r.m.CityID
}

func (r *debitCardInfoResolver) Branch() string {
	return r.m.Branch
}
