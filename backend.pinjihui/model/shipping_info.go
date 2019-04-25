package model

type ShippingInfo struct {
	ID              string
	Order_id        string
	Company         *string
	Delivery_number *string
	Images          *string
}

type ShippingInfoARR struct {
	ShippingInfo
	Images *[]string
}
