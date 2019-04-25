package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type SpikeRepository struct {
	BaseRepository
}

func NewSpikeRepository(db *sqlx.DB, log *logging.Logger) *SpikeRepository {
	return &SpikeRepository{BaseRepository{db: db, log: log}}
}

// 设置秒杀
func (r *SpikeRepository) SaveSpike(ctx context.Context, sp *model.Spike) (*model.Spike, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	productType, _, err := L("product").(*ProductRepository).FindProductType(ctx, sp.Product_id)
	if err != nil {
		return nil, err
	}

	if *productType != "simple" {
		return nil, errors.New("秒杀活动只能添加子商品")
	}

	var SQL string
	if sp.ID == "" {
		sp.ID = xid.New().String()
		if *usertype == "provider" || *usertype == "agent" {
			return nil, errors.New("您不能设置秒杀活动")
		}
		if *usertype == "ally" {
			sp.Merchant_id = *gc.CurrentUser(ctx)
		}
		if *usertype == "admin" {
			merchantID, err := r.GetMerchantID(sp.Product_id)
			if err != nil {
				return nil, err
			}
			sp.Merchant_id = *merchantID
		}
		expire, err := r.FindExpire(sp.Product_id, sp.Merchant_id)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			SQL = util.InsertSQLBuild(sp, "spikes", []string{})
		} else {
			t1, err := time.Parse(time.RFC3339, *expire)
			if err != nil {
				return nil, err
			}
			t2 := t1.Unix()
			t3 := time.Now().Unix()
			if t2 >= t3 {
				return nil, errors.New("您选的商品正处于秒杀中，不能继续添加")
			} else {
				SQL = util.InsertSQLBuild(sp, "spikes", []string{})
			}
		}

	} else {
		merchantID, err := r.GetMerchantID(sp.Product_id)
		if err != nil {
			return nil, err
		}
		sp.Merchant_id = *merchantID
		SQL = util.UpdateSQLBuild(sp, "spikes", []string{})
	}
	if _, err := r.db.NamedExec(SQL, sp); err != nil {
		return nil, err
	}
	spike, err := r.FindByID(ctx, sp.ID)
	if err != nil {
		return nil, err
	}
	return spike, nil
}

func (r *SpikeRepository) FindByID(ctx context.Context, ID string) (*model.Spike, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	spike := &model.Spike{}
	if *usertype == "admin" {
		SQL = `SELECT * FROM spikes WHERE id = $1`
	} else {
		SQL = `SELECT * FROM spikes WHERE id = $1` + gc.WhereMerchant(ctx)
	}
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(spike); err != nil {
		return nil, err
	}
	return spike, nil
}

func (r *SpikeRepository) FindAll(ctx context.Context) ([]*model.Spike, error) {
	var SQL string
	spikes := make([]*model.Spike, 0)
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `SELECT sp.* FROM spikes sp left join users u on sp.merchant_id=u.id where u.type='provider' ORDER BY sp.expired_at ASC ;`
		if err := r.db.Select(&spikes, SQL); err != nil {
			return nil, err
		}
		return spikes, nil
	} else if *usertype == "ally" {
		SQL = `SELECT * FROM spikes where merchant_id=$1 ORDER BY expired_at ASC ;`
		if err := r.db.Select(&spikes, SQL, *gc.CurrentUser(ctx)); err != nil {
			return nil, err
		}
		return spikes, nil
	} else {
		return nil, errors.New("您不能执行此操作")
	}
}

// 根据商品ID获取上传商家ID
func (r *SpikeRepository) GetMerchantID(productID string) (*string, error) {
	var merchantID string
	SQL := `select merchant_id from products where id=$1`
	err := r.db.Get(&merchantID, SQL, productID)
	if err != nil {
		return nil, err
	}
	return &merchantID, nil
}

// 根据商品ID和merchantID找到过期时间
func (r *SpikeRepository) FindExpire(productID string, merchantID string) (*string, error) {
	var expire string
	SQL := `SELECT expired_at FROM spikes WHERE product_id = $1 and merchant_id = $2`
	err := r.db.Get(&expire, SQL, productID, merchantID)
	if err != nil {
		return nil, err
	}
	return &expire, nil
}
