package resolver

import (
    //"pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
    "time"
    "pinjihui.com/pinjihui/loader"
    "golang.org/x/net/context"
)

type spikesConnectionResolver struct {
    list []*model.Spike
    Connection
}

func (r *spikesConnectionResolver) Spikes() *[]*spikeResolver {
    res := make([]*spikeResolver, len(r.list))
    for i := range res {
        res[i] = &spikeResolver{r.list[i]}
    }
    return &res
}

type spikeResolver struct {
    m *model.Spike
}

func (s *spikeResolver) ID() graphql.ID {
    return graphql.ID(s.m.ID)
}

func (s *spikeResolver) Product(ctx context.Context) (*productResolver, error) {
    product, err := loader.LoadProduct(ctx, s.m.ProductID, s.m.MerchantID)
    if err != nil {
        return nil, err
    }
    return &productResolver{product}, nil
}
func (s *spikeResolver) Price() float64 {
    return s.m.Price
}
func (c *spikeResolver) StartAt() (graphql.Time, error) {
    res, err := time.Parse(time.RFC3339, c.m.StartAt)
    return graphql.Time{Time: res}, err
}

func (c *spikeResolver) ExpiredAt() (graphql.Time, error) {
    res, err := time.Parse(time.RFC3339, c.m.ExpiredAt)
    return graphql.Time{Time: res}, err
}

func (c *spikeResolver) IsEmpty() bool {
    return c.m.TotalCount == 0
}

func (c *spikeResolver) Count() int32 {
    return c.m.TotalCount
}
