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

type AttributeRepository struct {
	BaseRepository
}

func NewAttributeRepository(db *sqlx.DB, log *logging.Logger) *AttributeRepository {
	return &AttributeRepository{BaseRepository{db: db, log: log}}
}

func (r *AttributeRepository) SaveAttribute(ctx context.Context, newAttr *model.Attribute) (*model.Attribute, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "ally" || *usertype == "agent" {
		return nil, errors.New("加盟商或者代理商不能添加或修改属性")
	}
	if newAttr.ID == "" {
		newAttr.ID = xid.New().String()
		if *usertype == "provider" {
			SQL = util.InsertSQLBuild(newAttr, "attributes", []string{"Required", "Searchable", "Enabled", "Deleted"})
		} else {
			SQL = util.InsertSQLBuild(newAttr, "attributes", []string{"Required", "Searchable", "Enabled", "Deleted", "Merchant_id"})
		}
	} else {
		if *usertype == "provider" {
			SQL = util.UpdateSQLBuild(newAttr, "attributes", []string{"Required", "Searchable", "Enabled", "Deleted", "Merchant_id"}) + gc.WhereMerchant(ctx)
		} else {
			SQL = util.UpdateSQLBuild(newAttr, "attributes", []string{"Required", "Searchable", "Enabled", "Deleted", "Merchant_id"}) + gc.WhereMerchantNULL()
		}
	}
	if _, err := r.db.NamedExec(SQL, newAttr); err != nil {
		return nil, err
	}
	attr, err := r.FindByID(ctx, newAttr.ID)
	if err != nil {
		return nil, err
	}
	return attr, nil
}

func (r *AttributeRepository) FindByID(ctx context.Context, ID string) (*model.Attribute, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	attribute := &model.Attribute{}
	if *usertype == "admin" {
		SQL = `SELECT * FROM attributes WHERE id = $1`
	} else {
		SQL = `SELECT * FROM attributes WHERE id = $1` + gc.WhereMerchant(ctx)
	}
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(attribute); err != nil {
		return nil, err
	}
	return attribute, nil
}

func (r *AttributeRepository) Deleted(ctx context.Context, ID string) (bool, error) {
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
		SQL = `update attributes set deleted=not deleted where id=$1` + gc.WhereMerchant(ctx)
	} else {
		SQL = `update attributes set deleted=not deleted where id=$1` + gc.WhereMerchantNULL()
	}
	if _, err := r.db.Exec(SQL, ID); err != nil {
		return false, err
	}
	return true, nil
}

func (r *AttributeRepository) FindAll(ctx context.Context) ([]*model.Attribute, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	attr := make([]*model.Attribute, 0)
	if *usertype == "ally" || *usertype == "agent" {
		return nil, errors.New("加盟商或者代理商不能执行此操作")
	} else if *usertype == "provider" {
		SQL = `SELECT * FROM attributes WHERE deleted=false AND enabled=true` + gc.WhereMerchantOR(ctx)

	} else {
		SQL = `SELECT * FROM attributes`
	}
	if err := r.db.Select(&attr, SQL); err != nil {
		return nil, err
	}
	return attr, nil
}

func (p *AttributeRepository) FindByIDs(ctx context.Context, ids *[]string) ([]*model.Attribute, error) {
	attrs := make([]*model.Attribute, 0)
	SQL := util.NewSQLBuilder(attrs).BuildQuery() + ` WHERE id in (?)`
	query, args, err := sqlx.In(SQL, *ids)
	if err != nil {
		return nil, err
	}
	query = p.db.Rebind(query)
	if err = p.db.Select(&attrs, query, args...); err != nil {
		return nil, err
	}
	return attrs, nil
}
