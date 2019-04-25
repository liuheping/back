package resolver

import (
	"strings"

	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type attributeResolver struct {
	m *model.Attribute
}

func (r *attributeResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *attributeResolver) Name() string {
	return r.m.Name
}

func (r *attributeResolver) Type() string {
	return r.m.Type
}

func (r *attributeResolver) Required() bool {
	return r.m.Required
}

func (r *attributeResolver) Enabled() bool {
	return r.m.Enabled
}

func (r *attributeResolver) Searchable() bool {
	return r.m.Searchable
}

func (r *attributeResolver) DefaultValue() *string {
	return r.m.Default_value
}

func (r *attributeResolver) Options() *[]string {
	if r.m.Options == nil {
		return nil
	}
	ra := util.Substr(*r.m.Options)
	res := strings.Split(ra, ",")
	for x, y := range res {
		res[x] = util.TrimQuotes(y)
	}
	return &res
}

func (r *attributeResolver) Deleted() bool {
	return r.m.Deleted
}

func (r *attributeResolver) Code() string {
	return r.m.Code
}
