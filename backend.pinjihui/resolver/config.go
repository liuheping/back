package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
)

type configResolver struct {
	m *model.Config
}

func (r *configResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *configResolver) Name() string {
	return r.m.Name
}

func (r *configResolver) Code() string {
	return r.m.Code
}

func (r *configResolver) Value() string {
	return r.m.Value
}

func (r *configResolver) Sortorder() *int32 {
	return r.m.Sort_order
}

func (r *configResolver) Description() *string {
	return r.m.Description
}
