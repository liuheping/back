package resolver

import (
	"strings"

	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

// 获取图片域名
func (r *Resolver) GetImageHost(ctx context.Context) string {
	return ctx.Value("config").(*gc.Config).QiniuCndHost
}

// 验证图片路径格式
func CompleteUrl(ctx context.Context, url string) string {
	if !strings.HasPrefix(url, "http") {
		url = ctx.Value("config").(*gc.Config).QiniuCndHost + url
		return url
	}
	return url
}

// 获取统计信息
func (r *Resolver) GetStatistic(ctx context.Context) (*statisticResolver, error) {
	statistic := &model.Statistic{}

	total, err := rp.L("public").(*rp.PublicRepository).GetOrderTotal(ctx)
	if err != nil {
		return nil, err
	}
	statistic.Total = total

	cost, err := rp.L("public").(*rp.PublicRepository).GetTotalCost(ctx)
	if err != nil {
		return nil, err
	}
	statistic.Cost = cost

	discount, err := rp.L("public").(*rp.PublicRepository).GetTotalOfferAmount(ctx)
	if err != nil {
		return nil, err
	}
	statistic.Discount = discount

	bonus, err := rp.L("public").(*rp.PublicRepository).GetTotalBonus(ctx)
	if err != nil {
		return nil, err
	}
	statistic.Bonus = bonus

	return &statisticResolver{statistic}, nil
}
