package model

type Coupon struct {
    ID          string
    Description string
    Value       float64
    StartAt     string   `db:"start_at"`
    ExpiredAt   string   `db:"expired_at"`
    LimitAmount *float64 `db:"limit_amount"`
}

const (
    ForSharer     = "for_sharer"
    ForBeSharer   = "for_be_sharer"
    ForFirstLogin = "for_first_login"
)
