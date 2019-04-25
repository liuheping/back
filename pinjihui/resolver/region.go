package resolver

import (
    "pinjihui.com/pinjihui/model"
)

type regionResolver struct {
    m *model.Region
}

func (r *regionResolver) ID() int32 {
    return r.m.ID
}

func (r *regionResolver) Name() string {
    return r.m.Name
}
