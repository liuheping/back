package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type OperationLogRepository struct {
	BaseRepository
}

func NewOperationLogRepository(db *sqlx.DB, log *logging.Logger) *OperationLogRepository {
	return &OperationLogRepository{BaseRepository{db: db, log: log}}
}

func (p *OperationLogRepository) FindByID(ctx context.Context, ID string) (*model.OperationLog, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	log := &model.OperationLog{}
	if *usertype == "provider" || *usertype == "ally" || *usertype == "agent" {
		SQL = `SELECT * FROM operation_logs WHERE id = $1` + ` AND user_id = '` + *gc.CurrentUser(ctx) + `'`
	} else {
		SQL = `SELECT * FROM operation_logs WHERE id = $1`
	}
	row := p.db.QueryRowx(SQL, ID)
	if err := row.StructScan(log); err != nil {
		return nil, err
	}
	return log, nil
}

func (p *OperationLogRepository) Search(ctx context.Context, first *int32, offset *int32, search *model.OperationLogSearchInput, sort *model.OperationLogSortInput) ([]*model.OperationLog, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	logs := make([]*model.OperationLog, 0)
	fetchSize := util.GetInt32(first, defaultListFetchSize)
	builder := util.NewSQLBuilder(logs).Table("operation_logs")
	p.searchWhere(search, builder)
	if offset != nil {
		if *usertype == "admin" {
			SQL = builder.BuildQuery() + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		} else {
			SQL = builder.BuildQuery() + ` AND user_id = '` + *gc.CurrentUser(ctx) + `'` + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		}
		if err := p.db.Select(&logs, SQL, fetchSize, *offset); err != nil {
			return nil, err
		}

		return logs, nil
	} else {
		if *usertype == "admin" {
			SQL = builder.BuildQuery() + p.buildSortSQL(sort) + ` LIMIT $1;`
		} else {
			SQL = builder.BuildQuery() + ` AND user_id = '` + *gc.CurrentUser(ctx) + `'` + p.buildSortSQL(sort) + ` LIMIT $1`
		}
		if err := p.db.Select(&logs, SQL, fetchSize); err != nil {
			return nil, err
		}
		return logs, nil
	}
}

func (p *OperationLogRepository) buildSortSQL(sort *model.OperationLogSortInput) (s string) {
	if sort == nil {
		return
	}
	s = " ORDER BY " + sort.OrderBy + " " + util.GetString(sort.Sort, "ASC") + " "
	return
}

func (p *OperationLogRepository) searchWhere(search *model.OperationLogSearchInput, builder *util.SQLBuilder) {
	builder.WhereStruct(search, true).WhereRow("true=true")
	if search != nil && search.Key != nil {
		builder.WhereRow(fmt.Sprintf("action ILIKE '%%%s%%'", *search.Key))
	}
}

func (p *OperationLogRepository) Count(ctx context.Context, c *model.OperationLogSearchInput) (int, error) {
	var count int
	var SQL string
	builder := util.NewSQLBuilder(nil)
	p.searchWhere(c, builder)
	where := builder.BuildWhere()
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return 0, err
	}
	if *status != "normal" {
		return 0, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `SELECT count(*) FROM operation_logs ` + where
	} else {
		SQL = `SELECT count(*) FROM operation_logs ` + where + ` AND user_id = '` + *gc.CurrentUser(ctx) + `'`
	}
	if err := p.db.Get(&count, SQL); err != nil {
		return 0, err
	}
	return count, nil
}
