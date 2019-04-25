package resolver

import (
	"strings"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

type brandResolver struct {
	m *model.Brand
}

func (r *brandResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *brandResolver) Name() string {
	return r.m.Name
}

func (r *brandResolver) Thumbnail(ctx context.Context) string {
	if r.m.Thumbnail == "" {
		return ""
	}
	// res := ctx.Value("config").(*gc.Config).QiniuCndHost + r.m.Thumbnail
	// return res
	return CompleteUrl(ctx, r.m.Thumbnail)
}

func (r *brandResolver) Description() *string {
	return r.m.Description
}

func (r *brandResolver) Enabled() bool {
	return r.m.Enabled
}

func (r *brandResolver) Deleted() bool {
	return r.m.Deleted
}

func (r *brandResolver) Sort_order() *int32 {
	return r.m.Sort_order
}

func (r *brandResolver) BrandType() string {
	return r.m.Brand_type
}

func (r *brandResolver) Created_at() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Created_at)
	return &graphql.Time{Time: res}, err
}

func (r *brandResolver) Updated_at() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Updated_at)
	return &graphql.Time{Time: res}, err
}

func (r *brandResolver) Machine_types() *[]string {
	if r.m.Machine_types == nil {
		return nil
	}
	ra := util.Substr(*r.m.Machine_types)
	res := strings.Split(ra, ",")
	for x, y := range res {
		res[x] = util.TrimQuotes(y)
	}
	return &res
}

func (r *brandResolver) Series(ctx context.Context) (*[]*brandSeriesResolver, error) {
	sers, err := rp.L("brandseries").(*rp.BrandSeriesRepository).FindByBrandID(ctx, r.m.ID)
	if err != nil {
		return nil, err
	}
	l := make([]*brandSeriesResolver, len(sers))
	for i := range l {
		l[i] = &brandSeriesResolver{(sers)[i]}
	}
	return &l, nil
}

func (r *brandResolver) Second_price_ratio() float64 {
	return r.m.Second_price_ratio
}

func (r *brandResolver) Retail_price_ratio() float64 {
	return r.m.Retail_price_ratio
}
