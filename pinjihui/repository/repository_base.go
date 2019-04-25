package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
)

type BaseRepository struct {
    db          *sqlx.DB
    log         *logging.Logger
}

func (r *BaseRepository) count(table string) (int, error) {
    var count int
    userSQL := `SELECT count(*) FROM ` + table
    err := r.db.Get(&count, userSQL)
    return count, err
}

func (r *BaseRepository) GetIntFormDB(sqlStr string, args ...interface{}) (int, error) {
    var v int
    err := r.db.Get(&v, sqlStr, args...)
    return v, err
}