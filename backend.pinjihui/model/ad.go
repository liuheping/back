package model

type Ad struct {
	ID          string
	Image       string
	Link        *string
	Merchant_id *string
	Position    string
	Sort        int32
	Is_show     bool
}
