package model

type ShippingAddress struct {
	ID         string
	UserId     string   `db:"user_id"`
	Mobile     string   `valid:"numeric,required,length(11|11)~mobile_invalid"`
	Consignee  string   `valid:"required"`
	Address    *Address `db:"*address"`
	Zipcode    *string
	IsDefault  bool `db:"is_default"`
	Created_at string
	Updated_at string
}

type ShippingAddressDB struct {
	Address string
	ShippingAddress
}
