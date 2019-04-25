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

type CouponRepository struct {
	BaseRepository
}

func NewCouponRepository(db *sqlx.DB, log *logging.Logger) *CouponRepository {
	return &CouponRepository{BaseRepository{db: db, log: log}}
}

func (r *CouponRepository) SaveCoupon(ctx context.Context, coupon *model.Coupon) (*model.Coupon, error) {
	var SQL string
	coupon.Merchant_id = gc.CurrentUser(ctx)
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if coupon.ID == "" {
		coupon.ID = xid.New().String()
		if *usertype == "admin" {
			// count, err := r.CountForNewer()
			// if err != nil {
			// 	return nil, err
			// }
			// if count >= 1 {
			// 	return nil, errors.New("只能创建一种新用户优惠券")
			// }
			SQL = util.InsertSQLBuild(coupon, "coupons", []string{"Created_at", "Updated_at", "Merchant_id"})
		} else {
			if coupon.Type == "for_newer" || coupon.Type == "for_inviter" {
				return nil, errors.New("您不能创建这种类型的优惠券")
			}
			SQL = util.InsertSQLBuild(coupon, "coupons", []string{"Created_at", "Updated_at"})
		}
	} else {
		if *usertype == "admin" {
			SQL = util.UpdateSQLBuild(coupon, "coupons", []string{"Created_at", "Updated_at", "Merchant_id"}) + gc.WhereMerchantNULL()
		} else {
			SQL = util.UpdateSQLBuild(coupon, "coupons", []string{"Created_at", "Updated_at", "Merchant_id"}) + gc.WhereMerchant(ctx)
		}

	}
	if _, err := r.db.NamedExec(SQL, coupon); err != nil {
		return nil, err
	}
	con, err := r.FindByID(ctx, coupon.ID)
	if err != nil {
		return nil, err
	}
	return con, nil
}

func (r *CouponRepository) FindByID(ctx context.Context, ID string) (*model.Coupon, error) {
	var SQL string
	coupon := &model.Coupon{}
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `SELECT * FROM coupons WHERE id = $1`
	} else {
		SQL = `SELECT * FROM coupons WHERE id = $1` + gc.WhereMerchant(ctx)
	}
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(coupon); err != nil {
		return nil, err
	}
	return coupon, nil
}

func (r *CouponRepository) FindAll(ctx context.Context) ([]*model.Coupon, error) {
	coupons := make([]*model.Coupon, 0)
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL := `SELECT * FROM coupons ORDER BY created_at DESC ;`
		err := r.db.Select(&coupons, SQL)
		if err != nil {
			return nil, err
		}
	} else {
		SQL := `SELECT * FROM coupons where merchant_id=$1 ORDER BY created_at DESC;`
		err := r.db.Select(&coupons, SQL, gc.CurrentUser(ctx))
		if err != nil {
			return nil, err
		}
	}
	return coupons, nil
}

func (r *CouponRepository) CountForNewer() (int, error) {
	var count int
	SQL := `SELECT count(*) FROM coupons where type='for_newer'`
	err := r.db.Get(&count, SQL)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CouponRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `DELETE FROM coupons WHERE id=$1` + gc.WhereMerchantNULL()
	} else {
		SQL = `DELETE FROM coupons WHERE id=$1` + gc.WhereMerchant(ctx)
	}
	result, err := r.db.Exec(SQL, ID)
	if err != nil {
		return false, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("删除优惠券失败, 检查ID为 %s 的优惠券是否存在", ID)
	}
	return true, nil
}
