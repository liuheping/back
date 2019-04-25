package model

type Spike struct {
    ID         string
    ProductID  string `db:"product_id"`
    Price      float64
    StartAt    string `db:"start_at"`
    ExpiredAt  string `db:"expired_at"`
    TotalCount int32  `db:"total_count"`
    MerchantID string `db:"merchant_id"`
}
