package resolver

import (
    "golang.org/x/net/context"
    gc "pinjihui.com/pinjihui/context"
    rp "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/service"
    "pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
)

func (r *Resolver) Brands(args struct {
    Type       string
    First      *int32
    MerchantID *string
    CateID     *string
}) (*[]*brandResolver, error) {
    brans, err := rp.L("brand").(*rp.BrandRepository).List(args.Type, args.First, args.MerchantID, args.CateID)
    if err != nil {
        return nil, err
    }
    bransr := make([]*brandResolver, len(brans))
    for i, v := range brans {
        bransr[i] = &brandResolver{v}
    }
    return &bransr, nil
}

func (r *Resolver) MachineTypes(args struct{ Series graphql.ID }) ([]string, error) {
    return rp.L("brand").(*rp.BrandRepository).FindMachineTypes(string(args.Series))
}

func (r *Resolver) Categories(args struct {
    MachineTypeSeries *string
    MerchantID        *string
    ParentID          *string
}) ([]*categoryResolver, error) {
    root, err := rp.L("category").(*rp.CategoryRepository).GetTree(args.MachineTypeSeries, args.MerchantID, args.ParentID)
    if err != nil {
        return nil, err
    }
    return GetChildResolvers(root)
}

func (r *Resolver) Cart(ctx context.Context) ([]*cartItemResolver, error) {
    return getCartItemResolver(ctx)
}

func getCartItemResolver(ctx context.Context) ([]*cartItemResolver, error) {
    items, err := rp.L("cart").(*rp.CartRepository).FindAll(ctx)
    if err != nil {
        return nil, err
    }
    itemResolvers := make([]*cartItemResolver, len(items))
    for i, v := range items {
        itemResolvers[i] = &cartItemResolver{v}
    }
    return itemResolvers, nil
}

func (r *Resolver) CartTotalCount(ctx context.Context) (int32, error) {
    c, err := rp.L("cart").(*rp.CartRepository).TotalCount(ctx)
    return int32(c), err
}

func (r *Resolver) Regions(args struct {
    Parent *int32
}) ([]*regionResolver, error) {
    regions, err := rp.L("region").(*rp.RegionRepository).FindByParentID(args.Parent)
    if err != nil {
        return nil, err
    }
    regionResolvers := make([]*regionResolver, len(regions))
    for i, v := range regions {
        regionResolvers[i] = &regionResolver{v}
    }
    return regionResolvers, nil
}

func (r *Resolver) QiniuUploadToken(ctx context.Context, args struct {
    Module string
    Ext    string
}) *QiniuUploadTokenResolver {
    gc.CheckAuth(ctx)
    token := service.GetQiniuUploadToken(ctx, args.Module, args.Ext)
    return &QiniuUploadTokenResolver{token}
}

type QiniuUploadTokenResolver struct {
    token *service.QiniuUploadToken
}

func (q *QiniuUploadTokenResolver) Token() string {
    return q.token.Token
}
func (q *QiniuUploadTokenResolver) Key() string {
    return q.token.Key
}

func (r *Resolver) Ads(args struct {
    Position   string
    MerchantID *string
}) ([]*adResolver, error) {
    ads, err := rp.L("ad").(*rp.ADRepository).List(args.Position, args.MerchantID)
    if err != nil {
        return nil, err
    }
    res := make([]*adResolver, len(ads))
    for i, v := range ads {
        res[i] = &adResolver{v}
    }
    return res, nil
}

type adResolver struct {
    m *model.AD
}

func (a *adResolver) Image(ctx context.Context) string {
    return completeUrl(ctx, a.m.Image)
}

func (a *adResolver) Link() *string {
    return a.m.Link
}
