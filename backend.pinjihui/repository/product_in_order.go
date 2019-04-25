package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"pinjihui.com/backend.pinjihui/model"
)

type ProductInOrderRepository struct {
	BaseRepository
}

func NewProductInOrderRepository(db *sqlx.DB, log *logging.Logger) *ProductInOrderRepository {
	return &ProductInOrderRepository{BaseRepository{db: db, log: log}}
}

//通过订单ID查找订单商品
func (p *ProductInOrderRepository) FindByOrderID(ID string) ([]*model.ProductInOrder, error) {
	products := []*model.ProductInOrder{}
	SQL := "SELECT * FROM order_products where order_id=$1"
	if err := p.db.Select(&products, SQL, ID); err != nil {
		return nil, err
	}
	return products, nil
}
