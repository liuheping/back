package resolver

import (
	"pinjihui.com/backend.pinjihui/model"
)

type incomeResolver struct {
	result interface{}
}

func (r *incomeResolver) ToProviderIncome() (*providerIncomeResolver, bool) {
	res, ok := r.result.(*providerIncomeResolver)
	return res, ok
}

func (r *incomeResolver) ToAllyIncome() (*allyIncomeResolver, bool) {
	res, ok := r.result.(*allyIncomeResolver)
	return res, ok
}

func (r *incomeResolver) ToAdminIncome() (*adminIncomeResolver, bool) {
	res, ok := r.result.(*adminIncomeResolver)
	return res, ok
}

type allyIncomeResolver struct {
	p *model.Order
}

func (r *allyIncomeResolver) Income() *float64 {
	return r.p.Ally_income
}

type providerIncomeResolver struct {
	p *model.Order
}

func (r *providerIncomeResolver) Income() *float64 {
	return r.p.Provider_income
}

type adminIncomeResolver struct {
	p *model.Order
}

func (r *adminIncomeResolver) ProviderIncome() *float64 {
	return r.p.Provider_income
}

func (r *adminIncomeResolver) AllyIncome() *float64 {
	return r.p.Ally_income
}
