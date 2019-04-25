package resolver

import (
    "pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
    "time"
)

type couponResolver struct {
    m *model.Coupon
}

func (c *couponResolver) ID() graphql.ID {
    return graphql.ID(c.m.ID)
}

func (c *couponResolver) Description() string {
    return c.m.Description
}

func (c *couponResolver) Value() float64 {
    return c.m.Value
}

func (c *couponResolver) LimitAmount() *float64 {
    return c.m.LimitAmount
}

func (c *couponResolver) StartAt() (string, error) {
    res, err := time.Parse(time.RFC3339, c.m.StartAt)
    return res.Format("2006.01.02"), err
}

func (c *couponResolver) ExpiredAt() (string, error) {
    res, err := time.Parse(time.RFC3339, c.m.ExpiredAt)
    return res.Format("2006.01.02"), err
}

type couponsConnectionResolver struct {
    list []*model.Coupon
    Connection
}

func (c *couponsConnectionResolver) Coupons() (*[]*couponResolver) {
    res := make([]*couponResolver, len(c.list))
    for i, v := range c.list {
        res[i] = &couponResolver{v}
    }
    return &res
}