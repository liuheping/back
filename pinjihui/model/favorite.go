package model

type Favorite struct {
    ID         string
    UserID     string  `db:"user_id"`
    MerchantID string  `db:"merchant_id"`
    ProductID  *string `db:"product_id"`
}
