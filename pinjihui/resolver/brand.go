package resolver

import (
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/repository"
    "sort"
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
    return completeUrl(ctx, r.m.Thumbnail)
}

func (r *brandResolver) Description() *string {
    return r.m.Description
}

func (r *brandResolver) MachineTypes() []string {
    return r.m.MachineTypes
}

func (r *brandResolver) MachineTypeSeries() ([]*machineTypeSeriesResolver, error) {
    series, err := repository.L("brand").(*repository.BrandRepository).FindSeries(r.m.ID)
    if err != nil {
        return nil, err
    }
    res := make([]*machineTypeSeriesResolver, len(series))
    for i, v := range series {
        res[i] = &machineTypeSeriesResolver{v}
    }
    return res, nil
}

type machineTypeSeriesResolver struct {
    series *model.BrandSeries
}

func (m *machineTypeSeriesResolver) ID() graphql.ID {
    return graphql.ID(m.series.ID)
}

func (m *machineTypeSeriesResolver) Name() string {
    return m.series.Name
}

func (m *machineTypeSeriesResolver) Image(ctx context.Context) *string {
    if m.series.Image == nil {
        return nil
    }
    res := completeUrl(ctx, *m.series.Image)
    return &res
}

func (m *machineTypeSeriesResolver) MachineTypes() []string {
    arr := m.series.GetMachineTypes()
    sort.Strings(arr)
    return arr
}
