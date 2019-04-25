package model

type OrderInput struct {
    ShippingAddrID string
    ItemIDs        []string
    SmOpt          []*ShppingMethodOption
    UsedCoupon     *[]string
    *InvOpt
}

type ShppingMethodOption struct {
    MerchantID     string
    ShippingMethod string
    Message        *string
}

type InvOpt struct {
    InvPayee    *string `db:"inv_payee"`
    InvType     string  `db:"inv_type"`
    InvTaxpayer *string `db:"inv_taxpayer_id"`
}
