package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type CommentRepository struct {
	BaseRepository
}

func NewCommentRepository(db *sqlx.DB, log *logging.Logger) *CommentRepository {
	return &CommentRepository{BaseRepository{db: db, log: log}}
}

func (p *CommentRepository) FindByID(ctx context.Context, ID string) (*model.Comment, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	comment := &model.Comment{}
	if *usertype == "provider" || *usertype == "ally" || *usertype == "agent" {
		SQL = `SELECT * FROM comments WHERE id = $1` + gc.WhereMerchant(ctx)
	} else {
		SQL = `SELECT * FROM comments WHERE id = $1`
	}
	row := p.db.QueryRowx(SQL, ID)
	if err := row.StructScan(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (p *CommentRepository) FindByProductID(ctx context.Context, ProductID string) ([]*model.Comment, error) {
	var SQL string
	comments := make([]*model.Comment, 0)
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `select * from comments where product_id = $1`
	} else {
		SQL = `select * from comments where product_id = $1` + gc.WhereMerchant(ctx)
	}
	if err := p.db.Select(&comments, SQL, ProductID); err != nil {
		return nil, err
	}
	return comments, nil
}

func (p *CommentRepository) Search(ctx context.Context, first *int32, offset *int32, search *model.CommentSearchInput, sort *model.CommentSortInput) ([]*model.Comment, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	comments := make([]*model.Comment, 0)
	fetchSize := util.GetInt32(first, defaultListFetchSize)
	builder := util.NewSQLBuilder(comments)
	p.searchWhere(search, builder)
	if offset != nil {
		if *usertype == "admin" {
			SQL = builder.BuildQuery() + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		} else {
			SQL = builder.BuildQuery() + gc.WhereMerchant(ctx) + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		}
		if err := p.db.Select(&comments, SQL, fetchSize, *offset); err != nil {
			return nil, err
		}
		return comments, nil
	} else {
		if *usertype == "admin" {
			SQL = builder.BuildQuery() + p.buildSortSQL(sort) + ` LIMIT $1;`
		} else {
			SQL = builder.BuildQuery() + gc.WhereMerchant(ctx) + p.buildSortSQL(sort) + ` LIMIT $1`
		}
		if err := p.db.Select(&comments, SQL, fetchSize); err != nil {
			return nil, err
		}
		return comments, nil
	}
}

func (p *CommentRepository) buildSortSQL(sort *model.CommentSortInput) (s string) {
	if sort == nil {
		return
	}
	s = " ORDER BY " + sort.OrderBy + " " + util.GetString(sort.Sort, "ASC") + " "
	return
}

func (p *CommentRepository) searchWhere(search *model.CommentSearchInput, builder *util.SQLBuilder) {
	//builder.WhereStruct(search, true)
	builder.WhereStruct(search, true).WhereRow("true=true")
	if search != nil && search.Key != nil {
		builder.WhereRow(fmt.Sprintf("content ILIKE '%%%s%%'", *search.Key))
	}
}

func (p *CommentRepository) Count(ctx context.Context, c *model.CommentSearchInput) (int, error) {
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
		SQL = `SELECT count(*) FROM comments ` + where
	} else {
		SQL = `SELECT count(*) FROM comments ` + where + gc.WhereMerchant(ctx)
	}
	if err := p.db.Get(&count, SQL); err != nil {
		return 0, err
	}
	return count, nil
}

//根据ID设置评论可见与否
func (r *CommentRepository) Visible(ctx context.Context, ID string) (bool, error) {
	SQL := `update comments set is_show= not is_show where id=$1` + gc.WhereMerchant(ctx)
	_, err := r.db.Exec(SQL, ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 回复评论
func (r CommentRepository) Reply(ctx context.Context, ID string, Content string) (*model.Comment, error) {
	SQL := `update comments set reply=$1,reply_time=$2 where id=$3` + gc.WhereMerchant(ctx)
	_, err := r.db.Exec(SQL, Content, time.Now(), ID)
	if err != nil {
		return nil, err
	}
	comment, err := r.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
