package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
)

type ADRepository struct {
    BaseRepository
}

func NewADRepository(db *sqlx.DB, log *logging.Logger) *ADRepository {
    return &ADRepository{BaseRepository{db: db, log: log}}
}

func (r *ADRepository) List(position string, merchantID *string) ([]*model.AD, error) {
    ads := make([]*model.AD, 0)
    sqls := `SELECT image, link FROM ads WHERE position=$1 AND merchant_id `
    args := []interface{}{position}
    if merchantID == nil {
        sqls += "IS NULL"
    } else {
        sqls += "=$2"
        args = append(args, merchantID)
    }
    sqls += " ORDER BY sort"
    err := r.db.Select(&ads, sqls, args...)
    return ads, err
}
