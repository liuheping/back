package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type AttributeSetRepository struct {
	BaseRepository
}

func NewAttributeSetRepository(db *sqlx.DB, log *logging.Logger) *AttributeSetRepository {
	return &AttributeSetRepository{BaseRepository{db: db, log: log}}
}

func (r *AttributeSetRepository) FindByID(ctx context.Context, ID string) (*model.AttributeSet, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	attributeset := &model.AttributeSet{}
	if *usertype == "admin" {
		SQL = `SELECT * FROM attribute_sets WHERE id = $1`
	} else {
		SQL = `SELECT * FROM attribute_sets WHERE id = $1` + gc.WhereMerchant(ctx)
	}
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(attributeset); err != nil {
		return nil, err
	}
	return attributeset, nil
}

func (r *AttributeSetRepository) FindAllAttributeBySetID(SetID string) ([]*model.Attribute, error) {
	attr := make([]*model.Attribute, 0)
	SQL := `SELECT attr.* FROM attribute_sets atts LEFT JOIN attributes attr on attr.id = any(atts.attribute_ids) WHERE atts.id=$1;`
	err := r.db.Select(&attr, SQL, SetID)
	if err != nil {
		return nil, err
	}
	return attr, nil
}

func (r *AttributeSetRepository) FindAll(ctx context.Context) ([]*model.AttributeSet, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	attr := make([]*model.AttributeSet, 0)
	if *usertype == "ally" || *usertype == "agent" {
		return nil, errors.New("您不能执行此操作")
	} else if *usertype == "provider" {
		SQL = `SELECT * FROM attribute_sets where deleted=false` + gc.WhereMerchantOR(ctx)
	} else {
		SQL = `SELECT * FROM attribute_sets`
	}
	if err := r.db.Select(&attr, SQL); err != nil {
		return nil, err
	}
	return attr, nil
}

// 更新或者添加属性集合
func (r *AttributeSetRepository) SaveAttributeSet(ctx context.Context, newAttrSet *model.AttributeSet) (*model.AttributeSet, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "ally" || *usertype == "agent" {
		return nil, errors.New("加盟商或者代理商不能添加或修改属性集合")
	}
	if newAttrSet.ID == "" {
		newAttrSet.ID = xid.New().String()
		if *usertype == "provider" {
			SQL = util.InsertSQLBuild(newAttrSet, "attribute_sets", []string{"Deleted"})
		} else {
			SQL = util.InsertSQLBuild(newAttrSet, "attribute_sets", []string{"Deleted", "Merchant_id"})
		}
	} else {
		if *usertype == "provider" {
			SQL = util.UpdateSQLBuild(newAttrSet, "attribute_sets", []string{"Merchant_id", "Deleted"}) + gc.WhereMerchant(ctx)
		} else {
			SQL = util.UpdateSQLBuild(newAttrSet, "attribute_sets", []string{"Merchant_id", "Deleted"}) //+ gc.WhereMerchantNULL()
		}
	}
	if _, err := r.db.NamedExec(SQL, newAttrSet); err != nil {
		return nil, err
	}
	attr, err := r.FindByID(ctx, newAttrSet.ID)
	if err != nil {
		return nil, err
	}
	return attr, nil
}

func (r *AttributeSetRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype == "ally" || *usertype == "agent" {
		return false, errors.New("加盟商或者代理商不能进行此操作")
	} else if *usertype == "provider" {
		SQL = `update attribute_sets set deleted=not deleted where id=$1` + gc.WhereMerchant(ctx)
	} else {
		SQL = `update attribute_sets set deleted=not deleted where id=$1` //+ gc.WhereMerchantNULL()
	}
	if _, err := r.db.Exec(SQL, ID); err != nil {
		return false, err
	}
	return true, nil
}
