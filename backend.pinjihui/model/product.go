package model

import graphql "github.com/graph-gophers/graphql-go"

type Product struct {
	ID               string
	Name             string
	Is_sale          bool
	Attribute_set_id *string
	Batch_price      *float64
	Second_price     *float64
	Category_id      string
	Related_ids      *string
	Content          *string
	Brand_id         string
	Deleted          bool
	Created_at       string
	Updated_at       string
	Tags             *string
	Attrs            *string
	Recommended      bool
	Sales_volume     int32
	Machine_types    *string
	Spec_1_name      *string
	Spec_2_name      *string
	Merchant_id      string
	Spec_1           *string
	Spec_2           *string
	Parent_id        *string
	Type             string
	Shipping_fee     *float64
}

type ProductSpecInput struct {
	Spec_1      string
	Spec_2      *string
	Batch_price float64
	Stock       int32
	Images      *[]string
}

type ProductARR struct {
	Product
	BatchPrice    *float64
	RelatedIds    *[]string
	Tags          *[]string
	Machine_types *[]string
	Images        *[]string
	Spec          *[]ProductSpecInput
	Stock         *int32
}

type ProductUpdateSpecInput struct {
	Product_id  *string
	Spec_1      string
	Spec_2      *string
	Batch_price float64
	Stock       int32
	Images      *[]string
}

type ProductUpdateARR struct {
	Product
	BatchPrice    *float64
	RelatedIds    *[]string
	Tags          *[]string
	Machine_types *[]string
	Images        *[]string
	Spec          *[]ProductUpdateSpecInput
	Stock         *int32
}

type ProductSearchInput struct {
	Key      *string     `db:"-"`
	Brand    *graphql.ID `db:"brand_id"`
	Category *graphql.ID `db:"category_id"`
}

type ProductSortInput struct {
	OrderBy string
	Sort    *string
}
