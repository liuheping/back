package model

type Coupon struct {
	ID            string
	Description   string
	Value         float64
	Created_at    string
	Updated_at    string
	Limit_amount  *float64
	Expired_at    *string
	Quantity      int32
	Type          string
	Start_at      *string
	Merchant_id   *string
	Validity_days *int32
}
