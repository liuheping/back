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

type BrandSeriesRepository struct {
	BaseRepository
}

func NewBrandSeriesRepository(db *sqlx.DB, log *logging.Logger) *BrandSeriesRepository {
	return &BrandSeriesRepository{BaseRepository{db: db, log: log}}
}

func (r *BrandSeriesRepository) SaveBrandSeries(ctx context.Context, series *model.BrandSeries) (*model.BrandSeries, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return nil, errors.New("只有管理员能创建更新品牌系列号")
	}
	if series.ID == "" {
		series.ID = xid.New().String()
		if series.Sort_order == nil {
			SQL = util.InsertSQLBuild(series, "brand_series", []string{"Sort_order"})
		}
		SQL = util.InsertSQLBuild(series, "brand_series", []string{})
	} else {
		SQL = util.UpdateSQLBuild(series, "brand_series", []string{})
	}

	// 更新或者删除系列号
	if _, err := r.db.NamedExec(SQL, series); err != nil {
		return nil, err
	}
	// 更新品牌表机型
	SQLStr := `UPDATE brands SET machine_types = array(SELECT unnest(machine_types) FROM brand_series WHERE brand_id = $1 ORDER BY sort_order) WHERE id = $2`
	if _, err := r.db.Exec(SQLStr, series.Brand_id, series.Brand_id); err != nil {
		return nil, err
	}

	ser, err := r.FindByID(ctx, series.ID)
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (r *BrandSeriesRepository) FindByID(ctx context.Context, ID string) (*model.BrandSeries, error) {
	series := &model.BrandSeries{}
	SQL := `SELECT * FROM brand_series WHERE id = $1`
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(series); err != nil {
		return nil, err
	}
	return series, nil
}

// 添加系列号时检查一个机型是否已经存在
func (r *BrandSeriesRepository) CheckMachineTypeForAdd(ctx context.Context, machinetype string) (bool, error) {
	var count int
	SQL := `SELECT count(*) FROM brand_series WHERE $1 = ANY(machine_types)`
	if err := r.db.Get(&count, SQL, machinetype); err != nil {
		return false, err
	}
	if count >= 1 {
		return false, errors.New("机型重复")
	}
	return true, nil
}

// 更新系列号时检查一个机型是否已经存在
func (r *BrandSeriesRepository) CheckMachineTypeForUpdate(ctx context.Context, machinetype string, seriesID string) (bool, error) {
	var count int
	SQL := `SELECT count(*) FROM brand_series WHERE $1 = ANY(machine_types) AND id != $2`
	if err := r.db.Get(&count, SQL, machinetype, seriesID); err != nil {
		return false, err
	}
	if count >= 1 {
		return false, errors.New("机型重复")
	}
	return true, nil
}

// 删除系列号
func (r *BrandSeriesRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return false, errors.New("只有管理员能删除系列号")
	}

	ser, err := r.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}

	SQL = `DELETE FROM brand_series WHERE id=$1`
	result, err := r.db.Exec(SQL, ID)
	if err != nil {
		return false, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("删除系列号失败，检查ID为 %s 的系列号是否存在", ID)
	}

	// 更新品牌表机型
	SQLStr := `UPDATE brands SET machine_types = array(SELECT unnest(machine_types) FROM brand_series WHERE brand_id = $1) WHERE id = $2`
	if _, err := r.db.Exec(SQLStr, ser.Brand_id, ser.Brand_id); err != nil {
		return false, err
	}

	return true, nil
}

// 根据品牌ID查询所有系列号
func (r *BrandSeriesRepository) FindByBrandID(ctx context.Context, brandID string) ([]*model.BrandSeries, error) {
	sers := make([]*model.BrandSeries, 0)
	SQL := `SELECT * FROM brand_series where brand_id = $1 ORDER BY sort_order;`
	err := r.db.Select(&sers, SQL, brandID)
	if err != nil {
		return nil, err
	}
	return sers, nil
}
