
package resolver

import (
    //"pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/model"
)

type attributeResolver struct {
    m *model.AttributeItem
}

func (r *attributeResolver) Name() string {
    return r.m.Name
}

func (r *attributeResolver) Value() *string {
    return r.m.Value
}
