package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
)

type AttributeRepository struct {
    BaseRepository
}

func NewAttributeRepository(db *sqlx.DB, log *logging.Logger) *AttributeRepository {
    return &AttributeRepository{BaseRepository{db: db, log: log}}
}

func (a *AttributeRepository) FindNamesByCodes(ids *[]string) ([]*model.Attribute, error) {
    names := make([]*model.Attribute, 0)
    sqls := `SELECT name, code FROM attributes WHERE code in (?) AND deleted=false AND enabled=true`
    query, args, err := sqlx.In(sqls, *ids)
    if err != nil {
        return nil, err
    }
    query = a.db.Rebind(query)
    if err = a.db.Select(&names, query, args...); err != nil {
        return nil, err
    }
    return names, nil
}
