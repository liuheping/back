package model

type CartItem struct {
    ID           string
    UserID       string `db:"user_id"`
    ProductID    string `db:"product_id"`
    ProductCount int32  `db:"product_count"`
    MerchantID   string `db:"merchant_id"`
    Price        float64
}
