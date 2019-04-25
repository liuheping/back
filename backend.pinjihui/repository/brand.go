package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type BrandRepository struct {
	BaseRepository
}

func NewBrandRepository(db *sqlx.DB, log *logging.Logger) *BrandRepository {
	return &BrandRepository{BaseRepository{db: db, log: log}}
}

func (r *BrandRepository) SaveBrand(ctx context.Context, brand *model.Brand) (*model.Brand, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return nil, errors.New("只有管理员能创建更新品牌")
	}

	if brand.Second_price_ratio < 0 || brand.Second_price_ratio > 1 {
		return nil, errors.New("二手价溢价率不在合法范围内")
	}

	if brand.Retail_price_ratio < 0 || brand.Retail_price_ratio > 1 {
		return nil, errors.New("售价溢价率不在合法范围内")
	}

	if brand.ID == "" {
		brand.ID = xid.New().String()
		SQL = util.InsertSQLBuild(brand, "brands", []string{"Deleted", "Enabled", "Created_at", "Updated_at", "Machine_types"})
	} else {
		SQL = util.UpdateSQLBuild(brand, "brands", []string{"Deleted", "Enabled", "Created_at", "Updated_at", "Machine_types"})
		// 修改了品牌就更新商品价格
		pSQL := `UPDATE products SET second_price = batch_price*$1 WHERE brand_id = $2`
		rSQL := `UPDATE rel_merchants_products SET retail_price = batch_price*$1 FROM products WHERE products.id=rel_merchants_products.product_id and products.brand_id = $2`
		x := float64(1) + brand.Second_price_ratio
		y := float64(1) + brand.Retail_price_ratio
		if _, err := r.db.Exec(pSQL, x, brand.ID); err != nil {
			return nil, errors.New("设置商品二手价失败")
		}
		if _, err := r.db.Exec(rSQL, y, brand.ID); err != nil {
			return nil, errors.New("设置商品售价失败")
		}
	}
	if _, err := r.db.NamedExec(SQL, brand); err != nil {
		return nil, err
	}
	br, err := r.FindByID(ctx, brand.ID)
	if err != nil {
		return nil, err
	}
	return br, nil
}

func (r *BrandRepository) FindByID(ctx context.Context, ID string) (*model.Brand, error) {
	brand := &model.Brand{}
	SQL := `SELECT * FROM brands WHERE id = $1`
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(brand); err != nil {
		return nil, err
	}
	return brand, nil
}

func (r *BrandRepository) FindAll(ctx context.Context) ([]*model.Brand, error) {
	brands := make([]*model.Brand, 0)
	SQL := `SELECT * FROM brands where deleted=false and enabled=true ORDER BY sort_order ASC ;`
	err := r.db.Select(&brands, SQL)
	if err != nil {
		return nil, err
	}
	return brands, nil
}

func (r *BrandRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return false, errors.New("只有管理员能删除品牌")
	}
	SQL := `update brands set deleted=not deleted where id=$1`
	if _, err := r.db.Exec(SQL, ID); err != nil {
		return false, nil
	}
	return true, nil
}

// 根据机器型号查找品牌
func (r *BrandRepository) FindByMachineType(ctx context.Context, machinetype string) (*model.Brand, error) {
	brand := &model.Brand{}
	SQL := `SELECT * FROM brands WHERE $1 = ANY(machine_types)`
	row := r.db.QueryRowx(SQL, machinetype)
	if err := row.StructScan(brand); err != nil {
		return nil, err
	}
	return brand, nil
}
