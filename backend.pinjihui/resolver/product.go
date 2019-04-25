package resolver

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

type productResolver struct {
	m *model.Product
}

func (r *productResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *productResolver) Name() *string {
	return &r.m.Name
}

func (r *productResolver) IsSale(ctx context.Context) (*bool, error) {
	isSale, err := rp.L("product").(*rp.ProductRepository).CatIsSale(ctx, r.m.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return isSale, nil
}

func (r *productResolver) AttributeSet(ctx context.Context) *attributeSetResolver {
	if r.m.Attribute_set_id == nil {
		return nil
	}
	attrset, err := rp.L("attributeset").(*rp.AttributeSetRepository).FindByID(ctx, *r.m.Attribute_set_id)
	if err != nil {
		return nil
	}
	return &attributeSetResolver{attrset}
}

func (r *productResolver) Category(ctx context.Context) (*categoryResolver, error) {
	cat, err := rp.L("category").(*rp.CategoryRepository).FindByID(ctx, r.m.Category_id)
	if err != nil {
		return nil, err
	}
	return &categoryResolver{cat}, nil
}

func (r *productResolver) RelatedIds() *[]graphql.ID {
	if r.m.Related_ids == nil {
		return nil
	}
	ra := util.Substr(*r.m.Related_ids)
	res := strings.Split(ra, ",")
	a := []graphql.ID{}
	for i := 0; i < len(res); i++ {
		a = append(a, graphql.ID(res[i]))
	}
	return &a
}

func (r *productResolver) Content() *string {
	return r.m.Content
}

func (r *productResolver) Brand(ctx context.Context) (*brandResolver, error) {
	brand, err := rp.L("brand").(*rp.BrandRepository).FindByID(ctx, r.m.Brand_id)
	if err != nil {
		return nil, err
	}
	return &brandResolver{brand}, nil
}

func (r *productResolver) Deleted() *bool {
	return &r.m.Deleted
}

func (r *productResolver) CreatedAt() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Created_at)
	return &graphql.Time{Time: res}, err
}

func (r *productResolver) UpdatedAt() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Updated_at)
	return &graphql.Time{Time: res}, err
}

func (r *productResolver) Tags() *[]string {
	if r.m.Tags == nil {
		return nil
	}
	ra := util.Substr(*r.m.Tags)
	res := strings.Split(ra, ",")
	for x, y := range res {
		res[x] = util.TrimQuotes(y)
	}
	return &res
}

func (r *productResolver) ProductImages() (*[]*productImageResolver, error) {
	images, err := rp.L("product").(*rp.ProductRepository).FindImagesById(r.m.ID)
	if err != nil {
		return nil, err
	}
	l := make([]*productImageResolver, len(images))
	for i := range l {
		l[i] = &productImageResolver{(images)[i]}
	}
	return &l, nil
}

func (r *productResolver) Stock(ctx context.Context) (*int32, error) {
	stock, err := rp.L("product").(*rp.ProductRepository).GetStock(ctx, r.m.ID)
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (r *productResolver) Sales_volume(ctx context.Context) (*int32, error) {
	SalesVolume, err := rp.L("product").(*rp.ProductRepository).GetSalesVolume(ctx, r.m.ID)
	if err != nil {
		return nil, err
	}
	return SalesVolume, nil
}

func (r *productResolver) Favorites(ctx context.Context) (*int32, error) {
	fav, err := rp.L("product").(*rp.ProductRepository).GetFavorites(ctx, r.m.ID)
	if err != nil {
		return nil, err
	}
	return fav, nil
}

func (r *productResolver) Attrs() *string {
	return r.m.Attrs
}

func (r *productResolver) Comments(ctx context.Context) (*[]*commentResolver, error) {
	comments, err := rp.L("comment").(*rp.CommentRepository).FindByProductID(ctx, r.m.ID)
	if err != nil {
		return nil, err
	}
	l := make([]*commentResolver, len(comments))
	for i := range l {
		l[i] = &commentResolver{(comments)[i]}
	}
	return &l, nil
}

func (r *productResolver) Specname1() *string {
	return r.m.Spec_1_name
}

func (r *productResolver) Specname2() *string {
	return r.m.Spec_2_name
}

func (r *productResolver) Spec1() *string {
	return r.m.Spec_1
}

func (r *productResolver) Spec2() *string {
	return r.m.Spec_2
}

func (r *productResolver) ParentID() *string {
	return r.m.Parent_id
}

func (r *productResolver) Type() string {
	return r.m.Type
}

func (r *productResolver) ShippingFee() *float64 {
	return r.m.Shipping_fee
}

func (r *productResolver) Price(ctx context.Context) (*priceResolver, error) {
	var a interface{}
	usertype, status, err := rp.L("public").(*rp.PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "ally" {
		a = &allyPriceResolver{r.m}
	}
	if *usertype == "provider" {
		a = &providerPriceResolver{r.m}
	}
	if *usertype == "admin" {
		a = &adminPriceResolver{r.m}
	}
	return &priceResolver{a}, nil
}

func (r *productResolver) Children() (*[]*productResolver, error) {
	products, err := rp.L("product").(*rp.ProductRepository).FindChildren(r.m.ID)
	if err != nil {
		return nil, err
	}
	l := make([]*productResolver, len(products))
	for i := range l {
		l[i] = &productResolver{(products)[i]}
	}
	return &l, nil
}

func (r *productResolver) Machine_types(ctx context.Context) (*[]*brandResolver, error) {
	if r.m.Machine_types == nil {
		return nil, nil
	}
	if *r.m.Machine_types == `{""}` {
		return nil, nil
	}
	ra := util.Substr(*r.m.Machine_types)
	res := strings.Split(ra, ",")
	// 去除引号
	for x, y := range res {
		res[x] = util.TrimQuotes(y)
	}

	brands := make([]*model.Brand, 0)
	for _, y := range res {
		brand, err := rp.L("brand").(*rp.BrandRepository).FindByMachineType(ctx, y)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		if len(brands) == 0 {
			*brand.Machine_types = "{" + y + "}"
			brands = append(brands, brand)
		} else {
			isexit := false
			for _, n := range brands {
				if brand.ID == n.ID {
					*n.Machine_types = "{" + util.Substr(*n.Machine_types) + "," + y + "}"
					isexit = true
					break
				}
			}
			if isexit == false {
				*brand.Machine_types = "{" + y + "}"
				brands = append(brands, brand)
			}
		}
	}

	l := make([]*brandResolver, len(brands))
	for i := range l {
		l[i] = &brandResolver{(brands)[i]}
	}
	return &l, nil
}
