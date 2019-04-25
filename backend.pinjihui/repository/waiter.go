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

type WaiterRepository struct {
	BaseRepository
}

func NewWaiterRepository(db *sqlx.DB, log *logging.Logger) *WaiterRepository {
	return &WaiterRepository{BaseRepository{db: db, log: log}}
}

// 通过ID查找客服
func (r *WaiterRepository) FindByID(ctx context.Context, ID string) (*model.Waiter, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	waiter := &model.Waiter{}
	if *usertype == "admin" {
		SQL = `SELECT * FROM waiters WHERE id = $1`
	} else {
		SQL = `SELECT * FROM waiters WHERE id = $1` + gc.WhereMerchant(ctx)
	}
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(waiter); err != nil {
		return nil, err
	}
	return waiter, nil
}

// 通过商家ID查找客服
func (r *WaiterRepository) FindByMerchantID(ctx context.Context) ([]*model.Waiter, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	waiters := make([]*model.Waiter, 0)
	if *usertype == "admin" {
		SQL = `SELECT * FROM waiters WHERE merchant_id IS NULL`
	} else {
		SQL = `SELECT * FROM waiters WHERE TRUE = TRUE` + gc.WhereMerchant(ctx)
	}
	if err := r.db.Select(&waiters, SQL); err != nil {
		return nil, err
	}
	return waiters, nil
}

func (r WaiterRepository) SaveWaiter(ctx context.Context, newWaiter *model.Waiter) (*model.Waiter, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	newWaiter.ID = xid.New().String()

	if *usertype == "admin" {
		SQL = util.InsertSQLBuild(newWaiter, "waiters", []string{"Merchant_id", "waiter_id", "Checked", "Deleted", "Remark"})
	} else {
		newWaiter.Merchant_id = gc.CurrentUser(ctx)
		SQL = util.InsertSQLBuild(newWaiter, "waiters", []string{"waiter_id", "Checked", "Deleted", "Remark"})
	}
	if _, err := r.db.NamedExec(SQL, newWaiter); err != nil {
		return nil, err
	}
	Waiter, err := r.FindByID(ctx, newWaiter.ID)
	if err != nil {
		return nil, err
	}
	return Waiter, nil
}

func (r WaiterRepository) UpdateWaiter(ctx context.Context, newWaiter *model.Waiter) (*model.Waiter, error) {
	var SQL string
	if newWaiter.ID == "" {
		return nil, errors.New("缺失修改对象ID")
	}
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	// 审核通过的客服不能修改
	Wa, err := r.FindByID(ctx, newWaiter.ID)
	if err != nil {
		return nil, err
	}
	if Wa.Checked == true {
		return nil, errors.New("审核通过的客服不能修改")
	}

	// 修改就需要重新审核
	newWaiter.Checked = false

	if *usertype == "admin" {
		SQL = util.UpdateSQLBuild(newWaiter, "waiters", []string{"Merchant_id", "waiter_id", "Deleted", "Remark"}) + gc.WhereMerchantNULL()
	} else {
		SQL = util.UpdateSQLBuild(newWaiter, "waiters", []string{"Merchant_id", "waiter_id", "Deleted", "Remark"}) + gc.WhereMerchant(ctx)
	}
	if _, err := r.db.NamedExec(SQL, newWaiter); err != nil {
		return nil, err
	}
	Waiter, err := r.FindByID(ctx, newWaiter.ID)
	if err != nil {
		return nil, err
	}
	return Waiter, nil
}

func (r WaiterRepository) Search(ctx context.Context, first *int32, offset *int32, search *model.WaiterSearchInput, sort *model.WaiterSortInput) ([]*model.Waiter, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	waiters := make([]*model.Waiter, 0)
	fetchSize := util.GetInt32(first, defaultListFetchSize)
	builder := util.NewSQLBuilder(waiters)
	r.searchWhere(search, builder)
	if offset != nil {
		if *usertype == "admin" {
			SQL = builder.BuildQuery() + r.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		} else {
			SQL = builder.BuildQuery() + gc.WhereMerchant(ctx) + r.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		}
		if err := r.db.Select(&waiters, SQL, fetchSize, *offset); err != nil {
			return nil, err
		}
		return waiters, nil
	} else {
		if *usertype == "admin" {
			SQL = builder.BuildQuery() + r.buildSortSQL(sort) + ` LIMIT $1;`
		} else {
			SQL = builder.BuildQuery() + gc.WhereMerchant(ctx) + r.buildSortSQL(sort) + ` LIMIT $1`
		}
		if err := r.db.Select(&waiters, SQL, fetchSize); err != nil {
			return nil, err
		}
		return waiters, nil
	}
}

func (r WaiterRepository) buildSortSQL(sort *model.WaiterSortInput) (s string) {
	if sort == nil {
		return
	}
	s = " ORDER BY " + sort.OrderBy + " " + util.GetString(sort.Sort, "ASC") + " "
	return
}

func (r WaiterRepository) searchWhere(search *model.WaiterSearchInput, builder *util.SQLBuilder) {
	//builder.WhereStruct(search, true)
	builder.WhereStruct(search, true).WhereRow("true=true")
	if search != nil && search.Key != nil {
		builder.WhereRow(fmt.Sprintf("remark ILIKE '%%%s%%'", *search.Key))
	}
}

func (r WaiterRepository) Count(ctx context.Context, c *model.WaiterSearchInput) (int, error) {
	var count int
	var SQL string
	builder := util.NewSQLBuilder(nil)
	r.searchWhere(c, builder)
	where := builder.BuildWhere()
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return 0, err
	}
	if *status != "normal" {
		return 0, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `SELECT count(*) FROM waiters ` + where
	} else {
		SQL = `SELECT count(*) FROM waiters ` + where + gc.WhereMerchant(ctx)
	}
	if err := r.db.Get(&count, SQL); err != nil {
		return 0, err
	}
	return count, nil
}

// 商家删除（注销）客服
func (r *WaiterRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}

	// 没有审核通过的直接物理删除,否则逻辑删除
	waiter, err := r.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}

	if *usertype == "admin" {
		if waiter.Checked == false {
			SQL = `DELETE FROM waiters WHERE id = $1` + gc.WhereMerchantNULL()
		} else {
			SQL = `UPDATE waiters SET deleted = true WHERE id = $1` + gc.WhereMerchantNULL()
		}
	} else {
		if waiter.Checked == false {
			SQL = `DELETE FROM waiters WHERE id = $1` + gc.WhereMerchant(ctx)
		} else {
			SQL = `UPDATE waiters SET deleted = true WHERE id = $1` + gc.WhereMerchant(ctx)
		}
	}

	if _, err := r.db.Exec(SQL, ID); err != nil {
		return false, err
	}

	return true, nil
}

// 审核客服，填入客服ID
func (r *WaiterRepository) CheckWaiter(ctx context.Context, ID string, waiter_id string, remark *string) (bool, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return false, errors.New("只有管理员能审核客服")
	}
	if waiter_id == "" {
		return false, errors.New("第三方的客服ID必须正确填写")
	}
	// 检查客服ID是否重复
	var count int
	countSQL := `SELECT count(*) FROM waiters where waiter_id = $1`
	if err := r.db.Get(&count, countSQL, waiter_id); err != nil {
		return false, err
	}
	if count > 0 {
		return false, errors.New("客服ID重复，请重新核实")
	}

	SQL = `UPDATE waiters SET checked = true , waiter_id = $1, remark = $2 WHERE id = $3`

	if _, err := r.db.Exec(SQL, waiter_id, remark, ID); err != nil {
		return false, err
	}
	return true, nil
}

// 删除第三方的客服帐号后，后台再删除数据
func (r *WaiterRepository) SaydeleteFromZhiMa(ctx context.Context, ID string) (bool, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}

	if *usertype != "admin" {
		return false, errors.New("只有管理员能进行此操作")
	}

	waiter, err := r.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}

	if waiter.Deleted == false {
		return false, errors.New("商家还没有删除客服，这里不能直接删除")
	}

	SQL = `DELETE FROM waiters WHERE id = $1 AND deleted = true AND checked = true`

	if _, err := r.db.Exec(SQL, ID); err != nil {
		return false, err
	}

	return true, nil
}
