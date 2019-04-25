package resolver

import (
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
    rp "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/util"
    "golang.org/x/net/context"
)

type orderProductItemResolver struct {
    m *model.OrderProduct
}

func (r *orderProductItemResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *orderProductItemResolver) ProductID() graphql.ID {
    return graphql.ID(r.m.ProductId)
}

func (r *orderProductItemResolver) Name() string {
    return r.m.ProductName
}

func (r *orderProductItemResolver) ProductCount() int32 {
    return r.m.ProductNumber
}

func (r *orderProductItemResolver) Price() string {
    return util.FmtMoney(r.m.ProductPrice)
}

func (r *orderProductItemResolver) Image(ctx context.Context) string {
    return getThumbnail(completeUrl(ctx, r.m.ProductImage))
}

func (r *orderProductItemResolver) ProductSpec() []*attributeResolver {
    attrs := make([]*attributeResolver, 0)
    if r.m.Spec1 == nil {
        return attrs
    }
    spec := &attributeResolver{&model.AttributeItem{*r.m.Spec1Name, r.m.Spec1}}
    attrs = append(attrs, spec)
    if r.m.Spec2Name != nil {
        spec = &attributeResolver{&model.AttributeItem{*r.m.Spec2Name, r.m.Spec2}}
        attrs = append(attrs, spec)
    }
    return attrs
}

func (r *orderProductItemResolver) IsCommented() (bool, error) {
    hasOne, err := rp.L("comment").(*rp.CommentRepository).HasOne(r.m.OrderId, r.m.ProductId)
    return hasOne, err
}

func (r *orderProductItemResolver) CommentID() *graphql.ID {
    if r.m.CommentID == nil {
        return nil
    }
    res := graphql.ID(*r.m.CommentID)
    return &res
}

func (r *orderProductItemResolver) ContainedShippingFee() bool {
    return r.m.ShippingFee != nil
}

type OrderProductItemConnnectionResolver struct {
    list []*model.OrderProduct
    Connection
}

func (op *OrderProductItemConnnectionResolver) Products() (*[]*orderProductItemResolver, error) {
    res := make([]*orderProductItemResolver, len(op.list))
    for i, v := range op.list {
        res[i] = &orderProductItemResolver{v}
    }
    return &res, nil
}
