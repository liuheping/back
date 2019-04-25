package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CreateBrand(ctx context.Context, args struct {
	NewBrand *model.Brand
}) (*brandResolver, error) {
	gc.CheckAuth(ctx)
	brand, err := rp.L("brand").(*rp.BrandRepository).SaveBrand(ctx, args.NewBrand)
	if err != nil {
		return nil, err
	}
	return &brandResolver{brand}, nil
}

func (r *Resolver) UpdateBrand(ctx context.Context, args struct {
	ID       graphql.ID
	NewBrand *model.Brand
}) (*brandResolver, error) {
	gc.CheckAuth(ctx)
	args.NewBrand.ID = string(args.ID)
	brand, err := rp.L("brand").(*rp.BrandRepository).SaveBrand(ctx, args.NewBrand)
	if err != nil {
		return nil, err
	}
	return &brandResolver{brand}, nil
}

func (r *Resolver) DeleteBrand(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("brand").(*rp.BrandRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}

// func (r *Resolver) CreateBrand(ctx context.Context, args struct {
// 	NewBrand *model.BrandARR
// }) (*brandResolver, error) {
// 	gc.CheckAuth(ctx)
// 	if args.NewBrand.Machine_types != nil {
// 		res := *args.NewBrand.Machine_types
// 		for _, y := range res {
// 			// 先判断输入机型数据库存在与否
// 			istrue, err := rp.L("brand").(*rp.BrandRepository).CheckMachineTypeForAdd(ctx, y)
// 			if istrue != true {
// 				return nil, err
// 			}
// 			// 判断一次性输入的机型有没有重复
// 			if util.RepeatCount(y, res) > 1 {
// 				return nil, errors.New("不能输入重复机型")
// 			}
// 		}
// 		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewBrand.Machine_types, "\",\""))
// 		args.NewBrand.Brand.Machine_types = &str
// 	}
// 	brand, err := rp.L("brand").(*rp.BrandRepository).SaveBrand(ctx, &args.NewBrand.Brand)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &brandResolver{brand}, nil
// }

// func (r *Resolver) UpdateBrand(ctx context.Context, args struct {
// 	ID       graphql.ID
// 	NewBrand *model.BrandARR
// }) (*brandResolver, error) {
// 	gc.CheckAuth(ctx)
// 	if args.NewBrand.Machine_types != nil {
// 		res := *args.NewBrand.Machine_types
// 		for _, y := range res {
// 			// 先判断输入机型数据库存在与否
// 			istrue, err := rp.L("brand").(*rp.BrandRepository).CheckMachineTypeForUpdate(ctx, y, string(args.ID))
// 			if istrue != true {
// 				return nil, err
// 			}
// 			// 判断一次性输入的机型有没有重复
// 			if util.RepeatCount(y, res) > 1 {
// 				return nil, errors.New("不能输入重复机型")
// 			}
// 		}
// 		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewBrand.Machine_types, "\",\""))
// 		args.NewBrand.Brand.Machine_types = &str
// 	}
// 	args.NewBrand.Brand.ID = string(args.ID)
// 	brand, err := rp.L("brand").(*rp.BrandRepository).SaveBrand(ctx, &args.NewBrand.Brand)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &brandResolver{brand}, nil
// }

// func (r *Resolver) DeleteBrand(ctx context.Context, args struct {
// 	ID graphql.ID
// }) (bool, error) {
// 	_, err := rp.L("brand").(*rp.BrandRepository).Deleted(ctx, string(args.ID))
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
