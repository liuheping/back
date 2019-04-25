package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

//获取某个商品的所有关联商家资料
func (r *Resolver) FindMerchantsByProductID(ctx context.Context, args struct{ ProductID string }) (*merchantProfileResolver, error) {
	merchants, err := rp.L("product").(*rp.ProductRepository).FindMerchants(args.ProductID)
	if err != nil {
		return nil, err
	}
	return &merchantProfileResolver{merchants}, nil
}

//获取某个商品的所有图片
func (r *Resolver) ProductImages(ctx context.Context, args struct{ ID string }) (*[]*productImageResolver, error) {
	images, err := rp.L("product").(*rp.ProductRepository).FindImagesById(args.ID)
	if err != nil {
		return nil, err
	}
	l := make([]*productImageResolver, len(images))
	for i := range l {
		l[i] = &productImageResolver{(images)[i]}
	}
	return &l, nil
}

//根据条件获取所有商品
func (r *Resolver) Products(ctx context.Context, args struct {
	First  *int32
	Offset *int32
	Search *model.ProductSearchInput
	Sort   *model.ProductSortInput
}) (*productsConnectionResolver, error) {
	fetchSize := util.GetInt32(args.First, DefaultPageSize)
	list, err := rp.L("product").(*rp.ProductRepository).Search(ctx, &fetchSize, args.Offset, args.Search, args.Sort)
	if err != nil {
		return nil, err
	}
	count, err := rp.L("product").(*rp.ProductRepository).Count(ctx, args.Search)
	if err != nil {
		return nil, err
	}
	var from *string
	var to *string
	if len(list) > 0 {
		from = &(list[0].ID)
		to = &(list[len(list)-1].ID)
	}
	return &productsConnectionResolver{m: list, totalCount: count, from: from, to: to, hasNext: util.If(len(list) == int(fetchSize), true, false).(bool)}, nil
}

//根据ID查找商品
func (r *Resolver) Product(ctx context.Context, args struct {
	ID string
}) (*productResolver, error) {
	product, err := rp.L("product").(*rp.ProductRepository).FindByID(args.ID)
	if err != nil {
		return nil, err
	}
	return &productResolver{product}, nil
}
