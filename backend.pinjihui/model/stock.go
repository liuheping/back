package model

import graphql "github.com/graph-gophers/graphql-go"

type Stock struct {
	Product_id   string
	Merchant_id  string
	Stock        int32
	Retail_price float64
	Origin_price *float64
	Created_at   string
	Sales_volume int32
	Is_sale      bool
	View_volume  int32
}

type StockSearchInput struct {
	Product_id  *graphql.ID
	Merchant_id *graphql.ID
}
