package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

type attributeSetResolver struct {
	m *model.AttributeSet
}

func (r *attributeSetResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *attributeSetResolver) Name() string {
	return r.m.Name
}

func (r *attributeSetResolver) Attributes() (*[]*attributeResolver, error) {
	return GetattributeResolver(r.m.ID)
}

func GetattributeResolver(ID string) (*[]*attributeResolver, error) {
	list, err := rp.L("attributeset").(*rp.AttributeSetRepository).FindAllAttributeBySetID(ID)
	if err != nil {
		return nil, err
	}
	listr := make([]*attributeResolver, len(list))
	for i, v := range list {
		listr[i] = &attributeResolver{v}
	}
	return &listr, nil
}
