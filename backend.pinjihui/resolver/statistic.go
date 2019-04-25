package resolver

import (
	"errors"

	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

type statisticResolver struct {
	s *model.Statistic
}

// 获取最畅销商品
func (r *statisticResolver) BestSaleProduct(ctx context.Context) (*[]*productResolver, error) {
	products, err := rp.L("public").(*rp.PublicRepository).BestSaleProduct(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*productResolver, len(products))
	for i := range l {
		l[i] = &productResolver{(products)[i]}
	}
	return &l, nil
}

// 获取收藏最多商品
func (r *statisticResolver) FavoriteProduct(ctx context.Context) (*[]*productResolver, error) {
	products, err := rp.L("public").(*rp.PublicRepository).FavoriteProduct(ctx)
	if err != nil {
		return nil, err
	}
	l := make([]*productResolver, len(products))
	for i := range l {
		l[i] = &productResolver{(products)[i]}
	}
	return &l, nil
}

func (r *statisticResolver) Total(ctx context.Context) *float64 {
	return r.s.Total
}

func (r *statisticResolver) Cost(ctx context.Context) *float64 {
	return r.s.Cost
}

func (r *statisticResolver) Discount(ctx context.Context) *float64 {
	return r.s.Discount
}

func (r *statisticResolver) Bonus(ctx context.Context) *float64 {
	return r.s.Bonus
}

func (r *statisticResolver) Profit(ctx context.Context) (*float64, error) {
	userType, status, err := rp.L("public").(*rp.PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *userType == "provider" {
		return nil, nil
	} else if *userType == "admin" {
		profit := *r.s.Total - *r.s.Cost - *r.s.Discount - *r.s.Bonus
		return &profit, nil
	} else {
		profit := *r.s.Total + *r.s.Bonus - *r.s.Cost - *r.s.Discount
		return &profit, nil
	}
}
