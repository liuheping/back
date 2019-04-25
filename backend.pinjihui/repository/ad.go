package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type AdRepository struct {
	BaseRepository
}

func NewAdRepository(db *sqlx.DB, log *logging.Logger) *AdRepository {
	return &AdRepository{BaseRepository{db: db, log: log}}
}

func (r *AdRepository) SaveAd(ctx context.Context, ad *model.Ad) (*model.Ad, error) {
	var SQL string
	ad.Merchant_id = gc.CurrentUser(ctx)
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if ad.ID == "" {
		ad.ID = xid.New().String()
		if *usertype == "admin" {
			SQL = util.InsertSQLBuild(ad, "ads", []string{"Merchant_id", "Is_show"})
		} else {
			SQL = util.InsertSQLBuild(ad, "ads", []string{"Is_show"})
		}
	} else {
		if *usertype == "admin" {
			SQL = util.UpdateSQLBuild(ad, "ads", []string{"Merchant_id", "Is_show"}) + gc.WhereMerchantNULL()
		} else {
			SQL = util.UpdateSQLBuild(ad, "ads", []string{"Merchant_id", "Is_show"}) + gc.WhereMerchant(ctx)
		}
	}
	result, err := r.db.NamedExec(SQL, ad)
	if err != nil {
		return nil, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return nil, err
	} else if af != 1 {
		return nil, fmt.Errorf("设置失败, 检查ID为 %s 的广告是否存在", ad.ID)
	}

	adv, err := r.FindByID(ctx, ad.ID)
	if err != nil {
		return nil, err
	}
	return adv, nil
}

func (r *AdRepository) FindByID(ctx context.Context, ID string) (*model.Ad, error) {
	var SQL string
	ad := &model.Ad{}
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `SELECT * FROM ads WHERE id = $1` + gc.WhereMerchantNULL()
	} else {
		SQL = `SELECT * FROM ads WHERE id = $1` + gc.WhereMerchant(ctx)
	}
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(ad); err != nil {
		return nil, err
	}
	return ad, nil
}

func (r *AdRepository) FindAll(ctx context.Context) ([]*model.Ad, error) {
	var SQL string
	ads := make([]*model.Ad, 0)
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `SELECT * FROM ads WHERE merchant_id IS NULL ORDER BY sort ASC `
	} else {
		SQL = `SELECT * FROM ads WHERE merchant_id = '` + *gc.CurrentUser(ctx) + `' ORDER BY sort ASC `
	}
	if err := r.db.Select(&ads, SQL); err != nil {
		return nil, err
	}
	return ads, nil
}

func (r *AdRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `DELETE FROM ads WHERE id=$1` + gc.WhereMerchantNULL()
	} else {
		SQL = `DELETE FROM ads WHERE id=$1` + gc.WhereMerchant(ctx)
	}
	result, err := r.db.Exec(SQL, ID)
	if err != nil {
		return false, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("删除失败, 检查ID为 %s 的广告是否存在", ID)
	}
	return true, nil
}

// 根据不同角色查找广告位置选项
func (r *AdRepository) FindAdPositionOptions(ctx context.Context) (*string, error) {
	var options string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		options = ctx.Value("config").(*gc.Config).AdminADPositionOptions
	} else {
		options = ctx.Value("config").(*gc.Config).MerchantADPositionOptions
	}
	return &options, nil
}

func (r *AdRepository) Isshow(ctx context.Context, ID string) (bool, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `UPDATE ads SET is_show = not is_show WHERE id=$1` + gc.WhereMerchantNULL()
	} else {
		SQL = `UPDATE ads SET is_show = not is_show WHERE id=$1` + gc.WhereMerchant(ctx)
	}
	result, err := r.db.Exec(SQL, ID)
	if err != nil {
		return false, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("设置失败, 检查ID为 %s 的广告是否存在", ID)
	}
	return true, nil
}
