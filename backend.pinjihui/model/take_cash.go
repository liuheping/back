package model

type TakeCash struct {
	UserID           string         `db:"user_id"`
	DebitCardInfo    *DebitCardInfo `db:"debit_card_info"`
	IsChecked        bool           `db:"is_checked"`
	DebitCardInfoRow string         `db:"debit_card_info"`
	Created_at       string
	Updated_at       string
}
