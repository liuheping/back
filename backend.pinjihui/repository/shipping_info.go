package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type ShippingInfoRepository struct {
	BaseRepository
}

func NewShippingInfoRepository(db *sqlx.DB, log *logging.Logger) *ShippingInfoRepository {
	return &ShippingInfoRepository{BaseRepository{db: db, log: log}}
}

// 更新或者创建物流信息
func (r *ShippingInfoRepository) SaveShippingInfo(ctx context.Context, shippingInfo *model.ShippingInfo) (*model.ShippingInfo, error) {
	var SQL string
	_, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if shippingInfo.ID == "" {
		shippingInfo.ID = xid.New().String()
		SQL = util.InsertSQLBuild(shippingInfo, "shipping_info", []string{})
	} else {
		SQL = util.UpdateSQLBuild(shippingInfo, "shipping_info", []string{})
	}
	if _, err := r.db.NamedExec(SQL, shippingInfo); err != nil {
		return nil, err
	}

	info, err := r.FindByID(ctx, shippingInfo.ID)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// 通过订单ID查找物流信息
func (r *ShippingInfoRepository) FindByOrderID(ctx context.Context, ID string) (*model.ShippingInfo, error) {
	info := &model.ShippingInfo{}
	SQL := `SELECT * FROM shipping_info WHERE order_id = $1`
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(info); err != nil {
		return nil, err
	}
	return info, nil
}

// 通过ID查找物流信息
func (r *ShippingInfoRepository) FindByID(ctx context.Context, ID string) (*model.ShippingInfo, error) {
	info := &model.ShippingInfo{}
	SQL := `SELECT * FROM shipping_info WHERE id = $1`
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(info); err != nil {
		return nil, err
	}
	return info, nil
}

// 通过订单ID删除物流信息
func (r *ShippingInfoRepository) DeletedByOrderID(ctx context.Context, OrderID string) (bool, error) {
	var SQL string
	_, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	SQL = `DELETE FROM shipping_info WHERE order_id=$1`
	result, err := r.db.Exec(SQL, OrderID)
	if err != nil {
		return false, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("删除物流信息失败,检查订单ID为 %s 的物流信息是否存在", OrderID)
	}
	return true, nil
}

// 通过ID删除物流信息
func (r *ShippingInfoRepository) DeletedByID(ctx context.Context, ID string) (bool, error) {
	var SQL string
	_, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	SQL = `DELETE FROM shipping_info WHERE id=$1`
	result, err := r.db.Exec(SQL, ID)
	if err != nil {
		return false, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("删除物流信息失败,检查ID为 %s 的物流信息是否存在", ID)
	}
	return true, nil
}
