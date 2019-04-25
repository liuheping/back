package model

type ProductImage struct {
    ID          string
    ProductId   string  `db:"product_id"`
    SmallImage  *string `db:"small_image"`
    MediumImage *string `db:"medium_image"`
    BigImage    string  `db:"big_image"`
}
