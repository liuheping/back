package model

type Category struct {
	ID        string
	ParentId  *string `db:"parent_id"`
	Name      string
	SortOrder *int32 `db:"sort_order"`
	Thumbnail *string
	Enabled   bool
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	Is_common bool
}
