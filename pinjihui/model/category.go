package model

type Category struct {
    ID          string
    ParentId    *string `db:"parent_id"`
    Name        string
    Thumbnail   *string
    Description *string
    Children    []*Category
    SortOrder   int     `db:"sort_order"`
    IsCommon    bool    `db:"is_common"`
}
