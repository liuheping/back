package model

type Spike struct {
	ID          string
	Product_id  string
	Price       float64
	Start_at    string
	Expired_at  string
	Total_count int32
	Merchant_id string
	Buy_limit   int32
}
