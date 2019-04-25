package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	logging "github.com/op/go-logging"
	"github.com/rs/xid"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
)

type PublicRepository struct {
	BaseRepository
}

func NewPublicRepository(db *sqlx.DB, log *logging.Logger) *PublicRepository {
	return &PublicRepository{BaseRepository{db: db, log: log}}
}

//获取当前用户类型
func (r *PublicRepository) GetCurrentUserType(ctx context.Context) (*string, *string, error) {
	var res struct {
		Type   string
		Status string
	}
	SQL := `select type, status from users where id=$1`
	row := r.db.QueryRowx(SQL, gc.CurrentUser(ctx))
	if err := row.StructScan(&res); err != nil {
		return nil, nil, err
	}
	return &res.Type, &res.Status, nil
}

// mutation 接口被访问就插入操作日志
func (r *PublicRepository) CreateOperationLog(ctx context.Context) {
	x := gc.QueryFileds(ctx)
	if len(x) == 0 {
		return
	}
	config := ctx.Value("config2").(*viper.Viper)
	for _, n := range x {
		z := config.Get("queryfiled." + n).(string)
		SQL := `insert into operation_logs(id,user_id,action) values ($1,$2,$3)`
		r.db.Exec(SQL, xid.New().String(), *gc.CurrentUser(ctx), z)
	}
}

// 最畅销商品
func (r *PublicRepository) BestSaleProduct(ctx context.Context) ([]*model.Product, error) {
	var SQL string
	userType, status, err := r.GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	if *userType == "admin" {
		SQL = `SELECT p.* FROM (SELECT product_id,SUM(sales_volume) sale FROM rel_merchants_products GROUP BY product_id ORDER BY SUM(sales_volume) DESC LIMIT 10 OFFSET 0) tt LEFT JOIN products p ON p."id" = tt.product_id ORDER BY tt.sale DESC`
	} else if *userType == "agent" {
		SQL = `SELECT p.* FROM rel_agents_products rmp LEFT JOIN products p ON rmp.product_id=p."id" WHERE rmp.agent_id = '` + *gc.CurrentUser(ctx) + `' ORDER BY rmp.sales_volume DESC LIMIT 10 OFFSET 0`
	} else {
		SQL = `SELECT p.* FROM rel_merchants_products rmp LEFT JOIN products p ON rmp.product_id=p."id" WHERE rmp.merchant_id = '` + *gc.CurrentUser(ctx) + `' ORDER BY rmp.sales_volume DESC LIMIT 10 OFFSET 0`
	}
	products := make([]*model.Product, 0)
	if err := r.db.Select(&products, SQL); err != nil {
		return nil, err
	}
	return products, nil
}

// 收藏最多商品
func (r *PublicRepository) FavoriteProduct(ctx context.Context) ([]*model.Product, error) {
	var SQL string
	userType, status, err := r.GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	if *userType == "admin" {
		SQL = `SELECT p.* FROM  (SELECT product_id,count(*) fav FROM favorites WHERE product_id is not null GROUP BY product_id ) tt LEFT JOIN products p on p.id = tt.product_id ORDER BY tt.fav DESC LIMIT 10 OFFSET 0`
	} else {
		SQL = `SELECT p.* FROM  (SELECT product_id,count(*) fav FROM favorites WHERE product_id is not null and merchant_id='` + *gc.CurrentUser(ctx) + `' GROUP BY product_id ) tt LEFT JOIN products p on p.id = tt.product_id ORDER BY tt.fav DESC LIMIT 10 OFFSET 0`
	}

	products := make([]*model.Product, 0)
	if err := r.db.Select(&products, SQL); err != nil {
		return nil, err
	}
	return products, nil
}

// 获取商家的订单总金额
func (r *PublicRepository) GetOrderTotal(ctx context.Context) (*float64, error) {
	var SQL, countSQL string
	var total float64
	var count int
	x := float64(0)

	userType, status, err := r.GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	if *userType == "agent" {
		return nil, nil
	}

	if *userType == "ally" {
		SQL = `SELECT sum(amount) FROM orders WHERE status='finish'` + gc.WhereMerchant(ctx)
		countSQL = `SELECT count(*) FROM orders WHERE status='finish'` + gc.WhereMerchant(ctx)
	} else if *userType == "provider" {
		SQL = `SELECT sum(provider_income) FROM orders WHERE status='finish' AND provider_income > 0`
		countSQL = `SELECT count(*) FROM orders WHERE status='finish' AND provider_income > 0`
	} else {
		SQL = `SELECT sum(amount) FROM orders WHERE status='finish' AND provider_income > 0`
		countSQL = `SELECT count(*) FROM orders WHERE status='finish' AND provider_income > 0`
	}

	if err := r.db.Get(&count, countSQL); err != nil {
		return &x, err
	}

	if count == 0 {
		return &x, nil
	} else {
		if err := r.db.Get(&total, SQL); err != nil {
			return &x, err
		}
		return &total, nil
	}
}

// 获取商家的订单总优惠
func (r *PublicRepository) GetTotalOfferAmount(ctx context.Context) (*float64, error) {
	var SQL, countSQL string
	var total float64
	var count int
	x := float64(0)

	userType, status, err := r.GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	if *userType == "agent" {
		return nil, nil
	}

	if *userType == "ally" {
		SQL = `SELECT sum(offer_amount) FROM orders WHERE status='finish'` + gc.WhereMerchant(ctx)
		countSQL = `SELECT count(*) FROM orders WHERE status='finish'` + gc.WhereMerchant(ctx)
	} else if *userType == "provider" {
		return nil, nil
	} else {
		SQL = `SELECT sum(offer_amount) FROM orders WHERE status='finish' AND provider_income > 0`
		countSQL = `SELECT count(*) FROM orders WHERE status='finish' AND provider_income > 0`
	}

	if err := r.db.Get(&count, countSQL); err != nil {
		return &x, err
	}

	if count == 0 {
		return &x, nil
	} else {
		if err := r.db.Get(&total, SQL); err != nil {
			return &x, err
		}
		return &total, nil
	}
}

// 获取商家的订单总成本
func (r *PublicRepository) GetTotalCost(ctx context.Context) (*float64, error) {
	var SQL, countSQL string
	var total float64
	var count int
	x := float64(0)

	userType, status, err := r.GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	if *userType == "agent" {
		return nil, nil
	}

	if *userType == "ally" {
		SQL = `SELECT SUM(op.second_price) FROM order_products op LEFT JOIN orders o ON o."id" = op.order_id LEFT JOIN users u ON u."id" = o.merchant_id WHERE o.status = 'finish' AND u."type" = 'ally' AND o.merchant_id = '` + *gc.CurrentUser(ctx) + `'`
		countSQL = `SELECT count(*) FROM order_products op LEFT JOIN orders o ON o."id" = op.order_id LEFT JOIN users u ON u."id" = o.merchant_id WHERE o.status = 'finish' AND u."type" = 'ally' AND o.merchant_id = '` + *gc.CurrentUser(ctx) + `'`
	} else if *userType == "provider" {
		return nil, nil
	} else {
		SQL = `SELECT sum(provider_income) FROM orders WHERE status='finish' AND provider_income > 0`
		countSQL = `SELECT count(*) FROM orders WHERE status='finish' AND provider_income > 0`
	}

	if err := r.db.Get(&count, countSQL); err != nil {
		return &x, err
	}

	if count == 0 {
		return &x, nil
	} else {
		if err := r.db.Get(&total, SQL); err != nil {
			return &x, err
		}
		return &total, nil
	}
}

// 获取商家的订单总提成
func (r *PublicRepository) GetTotalBonus(ctx context.Context) (*float64, error) {
	var SQL, countSQL string
	var total float64
	var count int
	x := float64(0)

	userType, status, err := r.GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	if *userType == "agent" {
		return nil, nil
	}

	if *userType == "ally" {
		SQL = `SELECT sum(ally_income) FROM orders WHERE status='finish'` + gc.WhereMerchant(ctx)
		countSQL = `SELECT count(*) FROM orders WHERE status='finish'` + gc.WhereMerchant(ctx)
	} else if *userType == "provider" {
		return nil, nil
	} else {
		SQL = `SELECT sum(ally_income) FROM orders WHERE status='finish'`
		countSQL = `SELECT count(*) FROM orders WHERE status='finish'`
	}

	if err := r.db.Get(&count, countSQL); err != nil {
		return &x, err
	}

	if count == 0 {
		return &x, nil
	} else {
		if err := r.db.Get(&total, SQL); err != nil {
			return &x, err
		}
		return &total, nil
	}
}
