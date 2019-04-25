package model

type AttributeSet struct {
	ID            string
	Name          string
	Merchant_id   *string
	Attribute_ids *string
	Deleted       bool
}

type AttributeSetARR struct {
	AttributeSet
	Attribute_ids *[]string
}
