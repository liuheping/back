package model

type Brand struct {
	ID                 string
	Name               string
	Thumbnail          string
	Description        *string
	Deleted            bool
	Enabled            bool
	Sort_order         *int32
	Created_at         string
	Updated_at         string
	Brand_type         string `db:"type"`
	Machine_types      *string
	Second_price_ratio float64
	Retail_price_ratio float64
}

type BrandARR struct {
	Brand
	Machine_types *[]string
}
