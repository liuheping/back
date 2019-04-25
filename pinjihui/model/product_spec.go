package model

type ProductSpec struct {
    ID    string
    Spec1 string  `db:"spec_1"`
    Spec2 *string `db:"spec_2"`
    Price float64 `db:"retail_price"`
    Stock int32
}
