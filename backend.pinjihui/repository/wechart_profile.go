package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
)

type WechartProfileRepository struct {
	BaseRepository
}

func NewWechartProfileRepository(db *sqlx.DB, log *logging.Logger) *WechartProfileRepository {
	return &WechartProfileRepository{BaseRepository{db: db, log: log}}
}

func (r *WechartProfileRepository) FindByUserID(ctx context.Context, ID string) (*model.WechartProfile, error) {
	wechart := &model.WechartProfile{}
	SQL := `SELECT * FROM wecharts WHERE user_id = $1`
	row := r.db.QueryRowx(SQL, ID)
	err := row.StructScan(wechart)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return wechart, nil
}
