package model

import graphql "github.com/graph-gophers/graphql-go"

type Comment struct {
	ID            string
	User_id       string
	Product_id    string
	Rank          *int32
	Shipping_rank *int32
	Service_rank  *int32
	Order_id      string
	Images        *string
	Content       string
	Is_show       bool
	Created_at    string
	User_ip       *string
	Reply         *string
	Reply_time    *string
	Merchant_id   string
	Updated_at    string
	Is_anonymous  bool
}

type CommentSearchInput struct {
	Key         *string     `db:"-"`
	User_id     *graphql.ID `db:"user_id"`
	Product_id  *graphql.ID `db:"product_id"`
	Order_id    *graphql.ID `db:"order_id"`
	Merchant_id *graphql.ID `db:"merchant_id"`
	Is_show     *bool       `db:"is_show"`
}

type CommentSortInput struct {
	OrderBy string
	Sort    *string
}
