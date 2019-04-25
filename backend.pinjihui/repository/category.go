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

type CategoryRepository struct {
	BaseRepository
}

func NewCategoryRepository(db *sqlx.DB, log *logging.Logger) *CategoryRepository {
	return &CategoryRepository{BaseRepository{db: db, log: log}}
}

func (r *CategoryRepository) List(parent_id *string) ([]*model.Category, error) {
	whereP := util.WhereNullable(parent_id)
	SQL := fmt.Sprintf(`SELECT id, parent_id, name, thumbnail,enabled,created_at,updated_at,is_common FROM %s WHERE enabled=true AND parent_id %s ORDER BY sort_order`,
		"categories", whereP)
	ms := []*model.Category{}
	if err := r.db.Select(&ms, SQL); err != nil {
		return nil, err
	}
	return ms, nil
}

func (r *CategoryRepository) SaveCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return nil, errors.New("只有管理员能创建更新分类")
	}
	if category.ID == "" {
		category.ID = xid.New().String()
		SQL = util.InsertSQLBuild(category, "categories", []string{"Sortorder", "Enabled", "CreatedAt", "UpdatedAt"})
	} else {
		SQL = util.UpdateSQLBuild(category, "categories", []string{"Sortorder", "Enabled", "CreatedAt", "UpdatedAt"})
	}
	if _, err := r.db.NamedExec(SQL, category); err != nil {
		return nil, err
	}
	cat, err := r.FindByID(ctx, category.ID)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, ID string) (*model.Category, error) {
	category := &model.Category{}
	SQL := `SELECT * FROM categories WHERE id = $1`
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return false, errors.New("只有管理员能删除分类")
	}
	SQL := `DELETE FROM categories WHERE id=$1`
	result, err := r.db.Exec(SQL, ID)
	if err != nil {
		return false, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("删除分类失败, 检查ID为 %s 的分类是否存在", ID)
	}
	return true, nil
}
