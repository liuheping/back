package resolver

import (
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type brandSeriesResolver struct {
	m *model.BrandSeries
}

func (r *brandSeriesResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *brandSeriesResolver) Brand_id() graphql.ID {
	return graphql.ID(r.m.Brand_id)
}

func (r *brandSeriesResolver) Series() string {
	return r.m.Series
}

func (r *brandSeriesResolver) Image(ctx context.Context) *string {
	if r.m.Image == nil {
		return nil
	}
	url := CompleteUrl(ctx, *r.m.Image)
	return &url
}

func (r *brandSeriesResolver) Machine_types() []string {
	ra := util.Substr(r.m.Machine_types)
	res := strings.Split(ra, ",")
	for x, y := range res {
		res[x] = util.TrimQuotes(y)
	}
	return res
}

func (r *brandSeriesResolver) Sort_order() int32 {
	return *r.m.Sort_order
}
