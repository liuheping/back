package resolver

import (
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
    "fmt"
    "pinjihui.com/pinjihui/repository"
    "golang.org/x/net/context"
    gc "pinjihui.com/pinjihui/context"
    "pinjihui.com/pinjihui/util"
)

type merchantResolver struct {
    m *model.Merchant
}

func (r *merchantResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *merchantResolver) CompanyName() string {
    if r.m.Type == model.PROVIDER {
        return model.Platform.CompanyName
    }
    return r.m.CompanyName
}

func (r *merchantResolver) CompanyAddress() string {
    if r.m.CompanyAddress == nil {
        return ""
    }
    return *r.m.CompanyAddress.RegionName + " " + r.m.CompanyAddress.Address
}

func (r *merchantResolver) DeliveryAddress() string {
    if r.m.DeliveryAddress == nil {
        return ""
    }
    return *r.m.DeliveryAddress.RegionName + " " + r.m.DeliveryAddress.Address
}

func (r *merchantResolver) CompanyImage(ctx context.Context) []string {
    urls := util.ParseArray(r.m.CompanyImage)
    for i, v := range urls {
        urls[i] = completeUrl(ctx, v)
    }
    return urls
}

func (r *merchantResolver) IsPlatform() bool {
    return r.m.IsPlatform()
}

func (r *merchantResolver) Logo(ctx context.Context) *string {
    if r.m.Logo == nil {
        return nil
    }
    url := completeUrl(ctx, *r.m.Logo)
    return &url
}

func (r *merchantResolver) Telephone() *string {
    return r.m.Telephone
}

func (r *merchantResolver) Distance(args struct{ Position *model.Location }) (string, error) {
    return GetDistance(r.m.ID, args.Position)
}

//todo 优化查询
func (r *merchantResolver) SalesVolume() (int32, error) {
    count, err := repository.L("merchant").(*repository.MerchantRepository).GetSalesVolume(r.m.ID)
    return int32(count), err
}

func (r *merchantResolver) CollectedQuantity() (int32, error) {
    count, err := repository.L("merchant").(*repository.MerchantRepository).GetCollectedQuantity(r.m.ID)
    return int32(count), err
}
func (r *merchantResolver) CommentQuantity() (int32, error) {
    count, err := repository.L("merchant").(*repository.MerchantRepository).GetCommentQuantity(r.m.ID)
    return int32(count), err
}

func GetDistance(id string, p *model.Location) (string, error) {
    if p == nil {
        return "", nil
    }
    d, err := repository.L("merchant").(*repository.MerchantRepository).GetDistance(id, p)
    return FmtDistance(d), err
}

func FmtDistance(d *float64) string {
    if d == nil {
        return ""
    }
    return fmt.Sprintf("%.2fkm", *d/1000)
}

func (r *merchantResolver) UserEdge(ctx context.Context) (*UserEdgeResolver, error) {
    var f *model.Favorite
    var err error
    if ctx.Value("is_authorized").(bool) {
        f, err = repository.L("favorite").(*repository.FavoriteRepository).FindByPMID(ctx, r.m.ID, nil)
        if err == gc.ErrNoRecord {
            err = nil
        }
    }
    return &UserEdgeResolver{f, nil}, nil
}

func (r *merchantResolver) Lat() *float64 {
    return r.m.Lat
}

func (r *merchantResolver) Lng() *float64 {
    return r.m.Lng
}

func (r *merchantResolver) Waiters() ([]string, error) {
    return repository.L("merchant").(*repository.MerchantRepository).FindWaitersByMerchantId(r.m.ID)
}

type UserEdgeResolver struct {
    f *model.Favorite
    p *model.Product
}

func (u *UserEdgeResolver) IsCollected(ctx context.Context) bool {
    return u.f != nil
}

func (u *UserEdgeResolver) FavoriteID() *graphql.ID {
    var id *graphql.ID
    if u.f != nil {
        res := graphql.ID(u.f.ID)
        id = &res
    }
    return id
}

func (u *UserEdgeResolver) IsPurchased(ctx context.Context) (bool, error) {
    if ctx.Value("is_authorized").(bool) {
        return repository.L("merchant").(*repository.MerchantRepository).IsPurchased(ctx, u.p.ID)
    }
    return false, nil
}