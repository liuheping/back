package resolver

import (
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
    rp "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/loader"
    "golang.org/x/net/context"
    "fmt"
    gc "pinjihui.com/pinjihui/context"
    "encoding/json"
    "strings"
    "pinjihui.com/pinjihui/util"
)

type productResolver struct {
    m *model.PaMCPair
}

func (r *productResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *productResolver) Name() string {
    return r.m.Name
}

func (r *productResolver) Price(ctx context.Context) float64 {
    return rp.L("product").(*rp.ProductRepository).GetPrice(r.m, ctx)
}

func (r *productResolver) RetailPrice(ctx context.Context) float64 {
    return util.GetFloat64(r.m.OriginPrice, r.m.Price)
}

func (r *productResolver) RelatedProducts() ([]*productResolver, error) {
    ids := r.m.GetRelatedIDArr()
    var paMCPairs []*model.PaMCPair
    if len(ids) > 0 {
        products, err := rp.L("product").(*rp.ProductRepository).FindByIDs(ids)

        if err != nil {
            return nil, err
        }
        paMCPairs = make([]*model.PaMCPair, len(products))
        for i, v := range products {
            productWithStock, err := rp.L("product").(*rp.ProductRepository).GetRelMerchantsProducts(v, r.m.MerchantID)
            if err != nil {
                return nil, err
            }
            paMCPairs[i] = productWithStock
        }
    }
    res := make([]*productResolver, len(paMCPairs))
    for i, v := range paMCPairs {
        res[i] = &productResolver{v}
    }
    return res, nil
}

func (r *productResolver) Content() *string {
    return r.m.Content
}

func (r *productResolver) Brand(ctx context.Context) (*brandResolver, error) {
    if r.m.BrandId == nil {
        return nil, nil
    }
    brand, err := loader.LoadBrand(ctx, *r.m.BrandId)
    if err != nil {
        return nil, err
    }
    return &brandResolver{brand}, nil
}

func (r *productResolver) Tags() []string {
    return r.m.GetTagArr()
}

func (r *productResolver) ProductImages() (*[]*productImageResolver, error) {
    images, err := rp.L("product").(*rp.ProductRepository).FindImagesById(r.m.ID)
    if err != nil {
        return nil, err
    }
    products := make([]*productImageResolver, len(images))
    for i := range images {
        products[i] = &productImageResolver{images[i]}
    }
    return &products, nil
}

func (r *productResolver) Attrs() (*[]*attributeResolver, error) {
    res := make([]*attributeResolver, 0)
    if r.m.Attrs == nil {
        return &res, nil
    }
    var attrsObj interface{}
    if err := json.Unmarshal([]byte(*r.m.Attrs), &attrsObj); err != nil {
        return nil, err
    }
    attrsMap := attrsObj.(map[string]interface{})
    keys := make([]string, 0, len(attrsMap))
    for k := range attrsMap {
        keys = append(keys, k)
    }
    if len(keys) == 0 {
        return &res, nil
    }
    names, err := rp.L("attr").(*rp.AttributeRepository).FindNamesByCodes(&keys)
    if err != nil {
        return nil, err
    }
    res = make([]*attributeResolver, len(names))
    var attr *model.AttributeItem
    var value string
    for i, v := range names {
        switch attrsMap[v.Code].(type) {
        case string:
            value = attrsMap[v.Code].(string)
        case []string:
            value = strings.Join(attrsMap[v.Code].([]string), ", ")
        }
        attr = &model.AttributeItem{v.Name, &value}
        res[i] = &attributeResolver{attr}
    }
    return &res, nil
}

func (r *productResolver) Spec() (*productSpecResolver, error) {
    if r.m.ParentID == nil {
        return nil, nil
    }
    items, err := rp.L("product").(*rp.ProductRepository).FindProductSpecs(r.m)
    if err != nil {
        return nil, err
    }
    if len(items) == 0 {
        return nil, nil
    }
    return &productSpecResolver{r.m, items}, nil
}

func (r *productResolver) SelfSpec() []*attributeResolver {
    spec := model.Spec{r.m.Spec1Name, r.m.Spec2Name, r.m.Spec1, r.m.Spec2}
    return resolverSpec(&spec)
}

func resolverSpec(product *model.Spec) []*attributeResolver {
    attrs := make([]*attributeResolver, 0)
    if product.Spec1Name == nil {
        return attrs
    }
    spec := &attributeResolver{&model.AttributeItem{*product.Spec1Name, product.Spec1}}
    attrs = append(attrs, spec)
    if product.Spec2Name != nil {
        spec = &attributeResolver{&model.AttributeItem{*product.Spec2Name, product.Spec2}}
        attrs = append(attrs, spec)
    }
    return attrs
}

func (r *productResolver) ClosestMerchant(ctx context.Context, args struct {
    Position    *model.Location
    ProductSpec *graphql.ID
}) (*edgeProductMerchantResolver, error) {
    if args.Position == nil {
        return nil, nil
    }
    ctx = context.WithValue(ctx, "position", args.Position)
    merchant, err := loader.LoadClosestMerchantsByPID(ctx, r.m.ID)
    fmt.Println(err)
    if err == gc.ErrNoRecord {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &edgeProductMerchantResolver{merchant}, nil
}

func (r *productResolver) Distance(args struct{ Position *model.Location }) (string, error) {
    return GetDistance(r.m.MerchantID, args.Position)
}

func (r *productResolver) RelMerchant(ctx context.Context) (*merchantResolver, error) {
    merchant, err := loader.LoadMerchant(ctx, r.m.MerchantID)
    if err != nil {
        return nil, err
    }
    return &merchantResolver{merchant}, nil
}

func (r *productResolver) Stock() int32 {
    return r.m.Stock
}

func (r *productResolver) ContainedShippingFee() bool {
    return r.m.ShippingFee != nil && *r.m.ShippingFee != 0
}

func (r *productResolver) OriginalPrice() string {
    if r.m.OriginPrice == nil {
        return util.FmtMoney(r.m.Price)
    }
    return util.FmtMoney(*r.m.OriginPrice)
}

func (r *productResolver) UserEdge(ctx context.Context) (*UserEdgeResolver, error) {
    var f *model.Favorite
    var err error
    if ctx.Value("is_authorized").(bool) {
        f, err = rp.L("favorite").(*rp.FavoriteRepository).FindByPMID(ctx, r.m.MerchantID, &r.m.ID)
        if err == gc.ErrNoRecord {
            err = nil
        }
    }
    return &UserEdgeResolver{f, &r.m.Product}, nil
}

func (r *productResolver) BestCommentRate() (string, error) {
    rank := int32(3)
    good, err := rp.L("comment").(*rp.CommentRepository).Count(r.m.ID, r.m.MerchantID, &rank)
    if err != nil {
        return "", err
    }
    all, err := rp.L("comment").(*rp.CommentRepository).Count(r.m.ID, r.m.MerchantID, nil)
    if err != nil {
        return "", err
    }
    if all == 0 {
        return "100%", nil
    }
    return fmt.Sprintf("%d%%", good/all*100), nil
}

func (r *productResolver) SalesVolume() int32 {
    return r.m.SalesVolume
}

func (r *productResolver) Spike() (*spikeResolver, error) {
    spike, err := rp.L("spike").(*rp.SpikeRepository).FindSpikeByPM(r.m.ID, r.m.MerchantID)
    if err == gc.ErrNoRecord {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &spikeResolver{spike}, nil
}

func (r *productResolver) ViewVolume() int32 {
    return r.m.ViewVolume
}
