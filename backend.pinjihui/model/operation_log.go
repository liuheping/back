package model

import graphql "github.com/graph-gophers/graphql-go"

type OperationLog struct {
	ID         string
	User_id    string
	Action     string
	Created_at string
}

type OperationLogSearchInput struct {
	Key     *string     `db:"-"`
	User_id *graphql.ID `db:"user_id"`
}

type OperationLogSortInput struct {
	OrderBy string
	Sort    *string
}
