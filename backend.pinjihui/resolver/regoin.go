package resolver

import (
	"pinjihui.com/backend.pinjihui/model"
)

type regoinResolver struct {
	m *model.Region
}

func (r *regoinResolver) ID() int32 {
	return r.m.ID
}

func (r *regoinResolver) Name() string {
	return r.m.Name
}

func (r *regoinResolver) Parent_id() int32 {
	return r.m.Parent_id
}

func (r *regoinResolver) Sort_order() *int32 {
	return &r.m.Sort_order
}
