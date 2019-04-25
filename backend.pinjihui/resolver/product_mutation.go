package resolver

import (
	"fmt"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

// 添加商品
func (r *Resolver) CreateProduct(ctx context.Context, args struct {
	NewPro *model.ProductARR
}) (*productResolver, error) {
	gc.CheckAuth(ctx)
	if args.NewPro.RelatedIds != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewPro.RelatedIds, "\",\""))
		args.NewPro.Product.Related_ids = &str
	}
	if args.NewPro.Tags != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewPro.Tags, "\",\""))
		args.NewPro.Product.Tags = &str
	}
	if args.NewPro.Machine_types != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewPro.Machine_types, "\",\""))
		args.NewPro.Product.Machine_types = &str

	}
	product, err := rp.L("product").(*rp.ProductRepository).Save(ctx, args.NewPro)
	if err != nil {
		return nil, err
	}
	return &productResolver{product}, nil
}

// 更新商品
func (r *Resolver) UpdateProduct(ctx context.Context, args struct {
	ID     graphql.ID
	NewPro *model.ProductUpdateARR
}) (*productResolver, error) {
	gc.CheckAuth(ctx)
	args.NewPro.Product.ID = string(args.ID)
	if args.NewPro.RelatedIds != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewPro.RelatedIds, "\",\""))
		args.NewPro.Product.Related_ids = &str
	}
	if args.NewPro.Tags != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewPro.Tags, "\",\""))
		args.NewPro.Product.Tags = &str
	}
	if args.NewPro.Machine_types != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.NewPro.Machine_types, "\",\""))
		args.NewPro.Product.Machine_types = &str

	}
	product, err := rp.L("product").(*rp.ProductRepository).Update(ctx, args.NewPro)
	if err != nil {
		return nil, err
	}
	return &productResolver{product}, nil
}

// 删除商品
func (r *Resolver) DeleteProduct(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("product").(*rp.ProductRepository).Deleted(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}

// 加盟商设置售价
func (r *Resolver) SetReatilPriceByAlly(ctx context.Context, args struct {
	ID    graphql.ID
	Price float64
}) (*productResolver, error) {
	product, err := rp.L("product").(*rp.ProductRepository).SetReatilPriceByAlly(ctx, args.Price, string(args.ID))
	if err != nil {
		return nil, err
	}
	return &productResolver{product}, nil
}

// 上下架
func (r *Resolver) SetProductIsSale(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := rp.L("product").(*rp.ProductRepository).SetIsSale(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
