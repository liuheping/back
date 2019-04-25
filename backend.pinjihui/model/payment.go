package model

type Payment struct {
	ID         string
	Pay_name   string
	Pay_code   string
	Pay_fee    *string
	Pay_desc   *string
	Sort_order *int32
	Enabled    bool
	Is_cod     bool
	Is_online  bool
	Deleted    bool
}
