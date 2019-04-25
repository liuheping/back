package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type OrderRepository struct {
	BaseRepository
}

func NewOrderRepository(db *sqlx.DB, log *logging.Logger) *OrderRepository {
	return &OrderRepository{BaseRepository{db: db, log: log}}
}

// 通过ID查找订单
func (r *OrderRepository) FindByID(ctx context.Context, ID string) (*model.Order, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	res := &model.OrderDB{}
	if *usertype == "admin" {
		SQL = `SELECT * FROM orders WHERE id = $1`
	} else {
		SQL = `SELECT * FROM orders WHERE id = $1` //+ gc.WhereMerchant(ctx) //这里如果加了的话在商家查询区域订单有问题
	}
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(res); err != nil {
		return nil, err
	}
	if res.Order.Address, err = model.NewOrderAddress(res.Address); err != nil {
		return nil, err
	}
	return &res.Order, nil
}

//查找子订单
func (p *OrderRepository) FindChildren(ctx context.Context, ID string) ([]*model.Order, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	orders := []*model.OrderDB{}
	if *usertype == "admin" {
		SQL = "SELECT * FROM orders WHERE parent_id = $1"
	} else {
		SQL = "SELECT * FROM orders WHERE parent_id = $1" + gc.WhereMerchant(ctx)
	}
	if err := p.db.Select(&orders, SQL, ID); err != nil {
		return nil, err
	}
	res := []*model.Order{}
	for _, y := range orders {
		var err error
		if y.Order.Address, err = model.NewOrderAddress(y.Address); err != nil {
			return nil, err
		}
		res = append(res, &y.Order)
	}
	return res, nil
}

//根据订单ID获取卖家资料
func (p *OrderRepository) FindMerchant(ID string) (*model.MerchantProfile, error) {
	profiles := &model.MerchantProfile{}
	SQL := `SELECT * FROM merchant_profiles mp WHERE mp.user_id=(SELECT merchant_id FROM orders WHERE id=$1)`
	row := p.db.QueryRowx(SQL, ID)
	err := row.StructScan(profiles)
	if err != nil {
		return nil, err
	}
	profiles.CompanyAddress, err = model.NewAddress(profiles.CompanyAddressRow)
	if err != nil {
		return nil, err
	}
	profiles.DeliveryAddress, err = model.NewAddress(profiles.DeliveryAddressRow)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

// 设置订单状态
func (p *OrderRepository) SetStatus(ctx context.Context, ID string, status string) (bool, error) {
	gc.CheckAuth(ctx)
	var SQL string
	order, err := p.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}
	if status == "shipped" {
		if order.Status != "paid" {
			return false, errors.New("还没付款，不能发货")
		}
		SQL = `update orders set status=$1,shipping_time=$2 where id=$3`
	} else {
		SQL = `update orders set status=$1,shipping_time=$2 where id=$3` + gc.WhereMerchant(ctx)
	}

	result, err := p.db.Exec(SQL, status, time.Now(), ID)
	if err != nil {
		return false, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("更新失败，检查订单是否存在, id: %s", ID)
	}
	return true, nil
}

// -------------------------------------------------

//根据条件查找订单
func (p *OrderRepository) Search(ctx context.Context, first *int32, offset *int32, search *model.OrderSearchInput, sort *model.OrderSortInput) ([]*model.Order, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	var SQL string
	orders := make([]*model.OrderDB, 0)
	res := []*model.Order{}
	fetchSize := util.GetInt32(first, defaultListFetchSize)
	builder := util.NewSQLBuilder(res)
	p.searchWhere(search, builder)
	if offset != nil {
		if *usertype == "admin" {
			SQL = builder.BuildQuery() + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		} else {
			SQL = builder.BuildQuery() + gc.WhereMerchant(ctx) + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		}
		if err := p.db.Select(&orders, SQL, fetchSize, *offset); err != nil {
			return nil, err
		}
		for _, y := range orders {
			var err error
			if y.Order.Address, err = model.NewOrderAddress(y.Address); err != nil {
				return nil, err
			}
			res = append(res, &y.Order)
		}
		return res, nil
	} else {
		if *usertype == "admin" {
			SQL = builder.BuildQuery() + p.buildSortSQL(sort) + ` LIMIT $1;`
		} else {
			SQL = builder.BuildQuery() + gc.WhereMerchant(ctx) + p.buildSortSQL(sort) + ` LIMIT $1;`
		}
		if err := p.db.Select(&orders, SQL, fetchSize); err != nil {
			return nil, err
		}
		for _, y := range orders {
			var err error
			if y.Order.Address, err = model.NewOrderAddress(y.Address); err != nil {
				return nil, err
			}
			res = append(res, &y.Order)
		}
		return res, nil
	}
}

func (p *OrderRepository) buildSortSQL(sort *model.OrderSortInput) (s string) {
	if sort == nil {
		return
	}
	s = " ORDER BY " + sort.OrderBy + " " + util.GetString(sort.Sort, "ASC") + " "
	return
}

func (p *OrderRepository) searchWhere(search *model.OrderSearchInput, builder *util.SQLBuilder) {
	builder.WhereStruct(search, true).WhereRow("merchant_id IS NOT NULL")
	if search != nil && search.Time != nil {
		builder.WhereRow(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", search.Time.Start, search.Time.End))
	}
}

// 根据不同角色和搜索条件获取订单数量
func (p *OrderRepository) Count(ctx context.Context, c *model.OrderSearchInput) (int, error) {
	var count int
	var SQL string
	builder := util.NewSQLBuilder(nil)
	p.searchWhere(c, builder)
	where := builder.BuildWhere()
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return 0, err
	}
	if *status != "normal" {
		return 0, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `SELECT count(*) FROM orders ` + where
	} else {
		SQL = `SELECT count(*) FROM orders ` + where + gc.WhereMerchant(ctx)
	}
	if err := p.db.Get(&count, SQL); err != nil {
		return 0, err
	}
	return count, nil
}

// 加盟商查看所在区域的相关订单（提成订单）
func (p *OrderRepository) SearchByArea(ctx context.Context, first *int32, offset *int32, search *model.OrderSearchInput, sort *model.OrderSortInput) ([]*model.Order, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "ally" {
		return nil, errors.New("您不能进行此操作")
	}
	var SQL string
	orders := make([]*model.OrderDB, 0)
	res := []*model.Order{}
	fetchSize := util.GetInt32(first, defaultListFetchSize)
	builder := util.NewSQLBuilder(nil)
	p.searchWhere(search, builder)
	where := builder.BuildWhere()
	if offset != nil {
		SQL = `SELECT * FROM (SELECT orders.* FROM orders LEFT JOIN users ON orders.user_id = users."id" LEFT JOIN users u ON orders.merchant_id = u."id"  WHERE users."type"='consumer' AND u."type" = 'provider' AND (orders.address).area_id = (SELECT (company_address).area_id FROM merchant_profiles WHERE user_id = $1)) AS tt` + where + p.buildSortSQL(sort) + ` LIMIT $2 OFFSET $3;`
		if err := p.db.Select(&orders, SQL, gc.CurrentUser(ctx), fetchSize, *offset); err != nil {
			return nil, err
		}
		for _, y := range orders {
			var err error
			if y.Order.Address, err = model.NewOrderAddress(y.Address); err != nil {
				return nil, err
			}
			res = append(res, &y.Order)
		}
		return res, nil
	} else {
		SQL = `SELECT * FROM (SELECT orders.* FROM orders LEFT JOIN users ON orders.user_id = users."id" LEFT JOIN users u ON orders.merchant_id = u."id"  WHERE users."type"='consumer' AND u."type" = 'provider' AND (orders.address).area_id = (SELECT (company_address).area_id FROM merchant_profiles WHERE user_id = $1)) AS tt` + where + p.buildSortSQL(sort) + ` LIMIT $2;`
		if err := p.db.Select(&orders, SQL, gc.CurrentUser(ctx), fetchSize); err != nil {
			return nil, err
		}
		for _, y := range orders {
			var err error
			if y.Order.Address, err = model.NewOrderAddress(y.Address); err != nil {
				return nil, err
			}
			res = append(res, &y.Order)
		}
		return res, nil
	}
}

// 根据不同角色和搜索条件获取所在区域订单数量
func (p *OrderRepository) CountByArea(ctx context.Context, c *model.OrderSearchInput) (int, error) {
	var count int
	var SQL string
	builder := util.NewSQLBuilder(nil)
	p.searchWhere(c, builder)
	where := builder.BuildWhere()
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return 0, err
	}
	if *status != "normal" {
		return 0, errors.New("用户状态不正常")
	}
	if *usertype != "ally" {
		return 0, errors.New("当前用户不允许执行此操作")
	}
	SQL = `SELECT COUNT(*) FROM (SELECT orders.* FROM orders LEFT JOIN users ON orders.user_id = users."id"	LEFT JOIN users u ON orders.merchant_id = u."id" WHERE users."type"='consumer' AND u."type" = 'provider' AND (orders.address).area_id = (SELECT (company_address).area_id FROM merchant_profiles WHERE user_id = $1)) AS tt` + where
	if err := p.db.Get(&count, SQL, gc.CurrentUser(ctx)); err != nil {
		return 0, err
	}
	return count, nil
}
