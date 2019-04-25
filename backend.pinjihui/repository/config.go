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

type ConfigRepository struct {
	BaseRepository
}

func NewConfigRepository(db *sqlx.DB, log *logging.Logger) *ConfigRepository {
	return &ConfigRepository{BaseRepository{db: db, log: log}}
}

func (r *ConfigRepository) SaveConfig(ctx context.Context, config *model.Config) (*model.Config, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return nil, errors.New("只有管理员能添加修改配置")
	}
	if config.ID == "" {
		config.ID = xid.New().String()
		SQL = util.InsertSQLBuild(config, "configs", []string{"Deleted"})
	} else {
		// if config.ID == "bc5mj9v2oau0r8ddf3q0" || config.ID == "bc70nn72oau4thff1evg" {
		// 	return nil, errors.New("敏感配置，不支持修改，如有需要，请联系开发人员")
		// }
		SQL = util.UpdateSQLBuild(config, "configs", []string{"Deleted"})
	}
	if _, err := r.db.NamedExec(SQL, config); err != nil {
		return nil, err
	}
	con, err := r.FindByID(ctx, config.ID)
	if err != nil {
		return nil, err
	}
	return con, nil
}

func (r *ConfigRepository) FindByID(ctx context.Context, ID string) (*model.Config, error) {
	config := &model.Config{}
	SQL := `SELECT * FROM configs WHERE id = $1`
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (r *ConfigRepository) FindByCode(ctx context.Context, code string) (*model.Config, error) {
	config := &model.Config{}
	SQL := `SELECT * FROM configs WHERE code = $1`
	row := r.db.QueryRowx(SQL, code)
	if err := row.StructScan(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (r *ConfigRepository) FindAll(ctx context.Context) ([]*model.Config, error) {
	configs := make([]*model.Config, 0)
	SQL := `SELECT * FROM configs where deleted=false ORDER BY sort_order ASC ;`
	if err := r.db.Select(&configs, SQL); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ConfigRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return false, errors.New("只有管理员能删除配置")
	}
	SQL := `update configs set deleted=not deleted where id=$1`
	if _, err := r.db.Exec(SQL, ID); err != nil {
		return false, err
	}
	return true, nil
}
