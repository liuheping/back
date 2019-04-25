package model

import graphql "github.com/graph-gophers/graphql-go"

type Waiter struct {
	ID          string
	Merchant_id *string
	Mobile      string
	Name        *string
	Waiter_id   *string
	Checked     bool
	Deleted     bool
	Remark      *string
}

type WaiterSearchInput struct {
	Key         *string     `db:"-"`
	Merchant_id *graphql.ID `db:"merchant_id"`
	Mobile      *string     `db:"mobile"`
	Name        *string     `db:"name"`
	Waiter_id   *graphql.ID `db:"product_id"`
	Checked     *bool       `db:"checked"`
	Deleted     *bool       `db:"deleted"`
}

type WaiterSortInput struct {
	OrderBy string
	Sort    *string
}
