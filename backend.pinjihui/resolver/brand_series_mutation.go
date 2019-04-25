package resolver

import (
	"errors"
	"fmt"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

func (r *Resolver) CreateBrandSeries(ctx context.Context, args struct {
	Series *model.BrandSeriesARR
}) (*brandSeriesResolver, error) {
	gc.CheckAuth(ctx)
	if args.Series.Machine_types != nil {
		res := args.Series.Machine_types
		for _, y := range res {
			// 先判断输入机型数据库存在与否
			istrue, err := rp.L("brandseries").(*rp.BrandSeriesRepository).CheckMachineTypeForAdd(ctx, y)
			if istrue != true {
				return nil, err
			}
			// 判断一次性输入的机型有没有重复
			if util.RepeatCount(y, res) > 1 {
				return nil, errors.New("不能输入重复机型")
			}
		}
		str := fmt.Sprintf("{\"%s\"}", strings.Join(args.Series.Machine_types, "\",\""))
		args.Series.BrandSeries.Machine_types = str
	}
	ser, err := rp.L("brandseries").(*rp.BrandSeriesRepository).SaveBrandSeries(ctx, &args.Series.BrandSeries)
	if err != nil {
		return nil, err
	}
	return &brandSeriesResolver{ser}, nil
}

func (r *Resolver) UpdateBrandSeries(ctx context.Context, args struct {
	ID     graphql.ID
	Series *model.BrandSeriesARR
}) (*brandSeriesResolver, error) {
	gc.CheckAuth(ctx)
	if args.Series.Machine_types != nil {
		res := args.Series.Machine_types
		for _, y := range res {
			// 先判断输入机型数据库存在与否
			istrue, err := rp.L("brandseries").(*rp.BrandSeriesRepository).CheckMachineTypeForUpdate(ctx, y, string(args.ID))
			if istrue != true {
				return nil, err
			}
			// 判断一次性输入的机型有没有重复
			if util.RepeatCount(y, res) > 1 {
				return nil, errors.New("不能输入重复机型")
			}
		}
		str := fmt.Sprintf("{\"%s\"}", strings.Join(args.Series.Machine_types, "\",\""))
		args.Series.BrandSeries.Machine_types = str
	}
	args.Series.BrandSeries.ID = string(args.ID)
	ser, err := rp.L("brandseries").(*rp.BrandSeriesRepository).SaveBrandSeries(ctx, &args.Series.BrandSeries)
	if err != nil {
		return nil, err
	}
	return &brandSeriesResolver{ser}, nil
}

func (r *Resolver) DeleteBrandSeries(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("brandseries").(*rp.BrandSeriesRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
