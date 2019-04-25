package model

import "pinjihui.com/pinjihui/util"

type Comment struct {
    ID           string
    UserID       string  `db:"user_id"`
    ProductID    string  `db:"product_id"`
    Rank         int32
    OrderID      string  `db:"order_id"`
    Content      string
    CreatedAt    string  `db:"created_at"`
    Reply        *string
    MerchantID   string  `db:"merchant_id"`
    ShippingRank int32   `db:"shipping_rank"`
    ServiceRank  int32   `db:"service_rank"`
    Images       *string
    UserIp       *string `fi:"-" db:"user_ip"`
    IsAnonymous  bool    `db:"is_anonymous" fi:"-"`
}

func (c *Comment) GetImageArr() []string {
    return util.ParseArray(c.Images)
}

type CommentInput struct {
    OrderProductID string `valid:"required"`
    Rank           int32  `valid:"range(1|3)"`
    Content        string `valid:"required"`
    ImageUrls      *[]string
    IsAnonymous    bool
    ShippingRank   int32  `valid:"range(1|5)"`
    ServiceRank    int32  `valid:"range(1|5)"`
}
