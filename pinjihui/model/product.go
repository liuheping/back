package model

import (
    "pinjihui.com/pinjihui/util"
)

type Product struct {
    ID             string
    Name           string
    AttributeSetId *string  `db:"attribute_set_id"`
    CategoryId     *string  `db:"category_id"`
    RelatedIds     *string  `db:"related_ids"`
    Content        *string
    BrandId        *string  `db:"brand_id"`
    Tags           *string
    Attrs          *string
    Spec1Name      *string  `db:"spec_1_name"`
    Spec2Name      *string  `db:"spec_2_name"`
    Spec1          *string  `db:"spec_1"`
    Spec2          *string  `db:"spec_2"`
    ParentID       *string  `db:"parent_id"`
    ShippingFee    *float64 `db:"shipping_fee"`
    SecondPrice    float64  `db:"second_price"`
    BatchPrice     float64  `db:"batch_price"`
}

type Spec struct {
    Spec1Name *string `db:"spec_1_name"`
    Spec2Name *string `db:"spec_2_name"`
    Spec1     *string `db:"spec_1"`
    Spec2     *string `db:"spec_2"`
}

func (p *Product) GetRelatedIDArr() []string {
    return util.ParseArray(p.RelatedIds)
}

func (p *Product) GetTagArr() []string {
    return util.ParseArray(p.Tags)
}

type PaMCPair struct {
    Product
    RelMerchantProduct
}

type RelMerchantProduct struct {
    Stock       int32
    Price       float64  `db:"retail_price"`
    MerchantID  string   `db:"merchant_id"`
    OriginPrice *float64 `db:"origin_price"`
    SalesVolume int32    `db:"sales_volume"`
    IsSale      bool     `db:"is_sale"`
    ViewVolume  int32    `db:"view_volume"`
}
