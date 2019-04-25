package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) Stocks(ctx context.Context, args struct{ Search *model.StockSearchInput }) (*[]*stockResolver, error) {
	stocks, err := rp.L("stock").(*rp.StockRepository).Find(args.Search)
	if err != nil {
		return nil, err
	}
	l := make([]*stockResolver, len(stocks))
	for i := range l {
		l[i] = &stockResolver{(stocks)[i]}
	}
	return &l, nil
}
