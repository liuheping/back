package model

type CashRequest struct {
	ID               string
	Merchant_id      string
	Amount           float64
	Status           string
	Reply            *string
	Note             *string
	Created_at       string
	Updated_at       string
	DebitCardInfo    *DebitCardInfo `db:"debit_card_info"`
	DebitCardInfoRow string         `db:"debit_card_info"`
}
