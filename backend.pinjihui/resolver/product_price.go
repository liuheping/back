package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

type priceResolver struct {
	result interface{}
}

func (r *priceResolver) ToProviderPrice() (*providerPriceResolver, bool) {
	res, ok := r.result.(*providerPriceResolver)
	return res, ok
}

func (r *priceResolver) ToAllyPrice() (*allyPriceResolver, bool) {
	res, ok := r.result.(*allyPriceResolver)
	return res, ok
}

func (r *priceResolver) ToAdminPrice() (*adminPriceResolver, bool) {
	res, ok := r.result.(*adminPriceResolver)
	return res, ok
}

type allyPriceResolver struct {
	p *model.Product
}

func (r *allyPriceResolver) SecondPrice() *float64 {
	return r.p.Second_price
}

func (r *allyPriceResolver) RetailPrice(ctx context.Context) (*float64, error) {
	price, err := rp.L("product").(*rp.ProductRepository).FindRetailPrice(ctx, r.p.ID)
	if err != nil {
		return nil, err
	}
	return price, nil
}

type providerPriceResolver struct {
	p *model.Product
}

func (r *providerPriceResolver) Price() *float64 {
	return r.p.Batch_price
}

type adminPriceResolver struct {
	p *model.Product
}

func (r *adminPriceResolver) RetailPrice(ctx context.Context) (*float64, error) {
	price, err := rp.L("product").(*rp.ProductRepository).FindRetailPriceForAdmin(ctx, r.p.ID)
	if err != nil {
		return nil, err
	}
	return price, nil
}

func (r *adminPriceResolver) SecondPrice() *float64 {
	return r.p.Second_price
}

func (r *adminPriceResolver) BatchPrice() *float64 {
	return r.p.Batch_price
}
