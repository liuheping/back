package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
)

type ShippingMethodRepository struct {
    BaseRepository
}

func NewShippingMethodRepository(db *sqlx.DB, log *logging.Logger) *ShippingMethodRepository {
    return &ShippingMethodRepository{BaseRepository{db: db, log: log}}
}

func (s *ShippingMethodRepository) FindAll() ([]*model.ShippingMethod, error) {
    sms := make([]*model.ShippingMethod, 0)
    sqls := `SELECT id, name, enable_for_platform FROM shipping_methods WHERE enabled=true`
    err := s.db.Select(&sms, sqls)
    return sms, err
}
