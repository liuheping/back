package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	//"pinjihui.com/backend.pinjihui/context"

	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type StockRepository struct {
	BaseRepository
}

func NewStockRepository(db *sqlx.DB, log *logging.Logger) *StockRepository {
	return &StockRepository{BaseRepository{db: db, log: log}}
}

//查找库存信息
func (r *StockRepository) Find(search *model.StockSearchInput) ([]*model.Stock, error) {
	stocks := make([]*model.Stock, 0)
	builder := util.NewSQLBuilder(stocks).Table("rel_merchants_products")
	builder.WhereStruct(search, true)
	SQL := builder.BuildQuery()
	err := r.db.Select(&stocks, SQL)
	if err != nil {
		return nil, err
	}
	return stocks, nil
}
