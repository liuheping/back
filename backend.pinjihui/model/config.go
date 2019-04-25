package model

type Config struct {
	ID          string
	Name        string
	Code        string
	Value       string
	Description *string
	Sort_order  *int32
	Deleted     bool
}
