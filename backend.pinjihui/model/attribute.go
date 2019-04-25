package model

type Attribute struct {
	ID            string
	Name          string
	Type          string
	Required      bool
	Default_value *string
	Options       *string
	Merchant_id   *string
	Searchable    bool
	Enabled       bool
	Deleted       bool
	Code          string
}

type AttributeARR struct {
	Attribute
	Options *[]string
}
