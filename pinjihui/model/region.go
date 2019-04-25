package model

type Region struct {
    ID          int32
    ParentId    int32   `db:"parent_id"`
    Name        string
    SortOrder   int32   `db:"sort_order"`
}
