package model

import graphql "github.com/graph-gophers/graphql-go"

type Order struct {
	ID      string
	User_id string
	// Order_status    string
	// Shipping_status string
	// Pay_status      string
	Postscript      *string
	Shipping_id     *string
	Shipping_name   *string
	Pay_code        string
	Pay_name        string
	Inv_payee       *string
	Inv_type        *string
	Amount          float64
	Shipping_fee    *float64
	Pay_fee         *float64
	Money_paid      *float64
	Order_amount    *float64
	Created_at      string
	Confirm_time    *string
	Pay_time        *string
	Shipping_time   *string
	Tax             *string
	Parent_id       *string
	Merchant_id     *string
	Address         *OrderAddress `db:"*address"`
	Note            *string
	Inv_taxpayer_id *string
	Inv_url         *string
	Used_coupon     *string
	Status          string
	Offer_amount    float64
	Provider_income *float64
	Ally_income     *float64
}

type OrderDB struct {
	Order
	Address string
}

type OrderSearchInput struct {
	// Key         *string     `db:"-"`
	User_id     *graphql.ID     `db:"user_id"`
	Shipping_id *graphql.ID     `db:"shipping_id"`
	Merchant_id *graphql.ID     `db:"merchant_id"`
	Pay_code    *string         `db:"pay_code"`
	Status      *string         `db:"status"`
	Time        *OrderTimeInput `db:"-"`
}

type OrderSortInput struct {
	OrderBy string
	Sort    *string
}

type OrderTimeInput struct {
	Start string
	End   string
}
