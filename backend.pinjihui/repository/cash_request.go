package repository

import (
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

type CashRequestRepository struct {
	BaseRepository
}

func NewCashRequestRepository(db *sqlx.DB, log *logging.Logger) *CashRequestRepository {
	return &CashRequestRepository{BaseRepository{db: db, log: log}}
}

// 提现申请
func (r *CashRequestRepository) SaveCashRequest(ctx context.Context, cr *model.CashRequest) (*model.CashRequest, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	// 只有10号、20号、30号能提现
	if time.Now().Day() != 8 || time.Now().Day() != 18 || time.Now().Day() != 28 {
		return nil, errors.New("只有8号、18号、28号能提现")
	}

	if *usertype != "ally" && *usertype != "provider" && *usertype != "agent" {
		return nil, errors.New("您不能提现")
	}

	if iscan := r.Can(ctx); iscan == false {
		return nil, errors.New("您有提现正在进行中，请等待结束后再操作")
	}
	cr.ID = xid.New().String()
	cr.Merchant_id = *gc.CurrentUser(ctx)
	cr.Status = "unchecked"
	balance, err := r.FindBalanceByMerchantID(*gc.CurrentUser(ctx))
	if err != nil {
		return nil, err
	}
	if cr.Amount < float64(0) {
		return nil, errors.New("提现金额不能小于0")
	}
	if cr.Amount > *balance {
		return nil, errors.New("余额不足")
	}
	SQL := util.InsertSQLBuild(cr, "cash_requests", []string{"Reply", "Created_at", "Updated_at", "DebitCardInfoRow"})
	if _, err := r.db.NamedExec(SQL, cr); err != nil {
		return nil, err
	}
	if _, err := r.ReduceBalanceByMerchantID(cr.Amount, cr.Merchant_id); err != nil {
		return nil, err
	}
	cashrequest, err := r.FindByID(ctx, cr.ID)
	if err != nil {
		return nil, err
	}
	return cashrequest, nil
}

func (r *CashRequestRepository) FindBalanceByMerchantID(ID string) (*float64, error) {
	var balance float64
	SQL := `SELECT balance FROM merchant_profiles WHERE user_id = $1`
	if err := r.db.Get(&balance, SQL, ID); err != nil {
		zero := float64(0)
		return &zero, err
	}
	return &balance, nil
}

func (r *CashRequestRepository) FindByID(ctx context.Context, ID string) (*model.CashRequest, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	res := &model.CashRequest{}
	if *usertype != "admin" {
		SQL = `SELECT * FROM cash_requests WHERE id = $1` + gc.WhereMerchant(ctx)
	} else {
		SQL = `SELECT * FROM cash_requests WHERE id = $1`
	}
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(res); err != nil {
		return nil, err
	}
	res.DebitCardInfo, err = model.NewDebitCardInfo(res.DebitCardInfoRow)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *CashRequestRepository) FindAll(ctx context.Context) ([]*model.CashRequest, error) {
	res := make([]*model.CashRequest, 0)
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL := `SELECT * FROM cash_requests ;`
		if err := r.db.Select(&res, SQL); err != nil {
			return nil, err
		}
	} else {
		SQL := `SELECT * FROM cash_requests where merchant_id=$1;`
		if err := r.db.Select(&res, SQL, gc.CurrentUser(ctx)); err != nil {
			return nil, err
		}
	}
	for _, v := range res {
		v.DebitCardInfo, err = model.NewDebitCardInfo(v.DebitCardInfoRow)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// 增加余额
func (r *CashRequestRepository) AddBalanceByMerchantID(amount float64, ID string) (bool, error) {
	SQL := `update merchant_profiles set balance=balance+$1 where user_id=$2`
	if _, err := r.db.Exec(SQL, amount, ID); err != nil {
		return false, err
	}
	return true, nil
}

// 减少余额
func (r *CashRequestRepository) ReduceBalanceByMerchantID(amount float64, ID string) (bool, error) {
	SQL := `update merchant_profiles set balance=balance-$1 where user_id=$2`
	if _, err := r.db.Exec(SQL, amount, ID); err != nil {
		return false, err
	}
	return true, nil
}

// 查找当前会话用户有无提现正在进行中,true ——> 没有进行中的 ——> 可以提现
func (r *CashRequestRepository) Can(ctx context.Context) bool {
	var count int
	SQL := `SELECT count(*) FROM cash_requests where merchant_id = $1 and status != $2 and status != $3`
	if err := r.db.Get(&count, SQL, gc.CurrentUser(ctx), "finished", "refused"); err != nil {
		return false
	}
	if count > 0 {
		return false
	} else {
		return true
	}
}

// 管理员拒绝提现
func (r *CashRequestRepository) SetRefused(ctx context.Context, ID string, Reply string) (bool, error) {
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
	cr, err := r.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}
	if cr.Status != "unchecked" {
		return false, errors.New("当前状态下不允许进行此操作")
	}

	SQL := `update cash_requests set status = $1, reply = $2 where id = $3`
	if _, err := r.db.Exec(SQL, "refused", Reply, ID); err != nil {
		return false, err
	}
	//拒绝后返还余额
	if _, err := r.AddBalanceByMerchantID(cr.Amount, cr.Merchant_id); err != nil {
		return false, err
	}
	return true, nil
}

// 管理员设置已打款
func (r *CashRequestRepository) SetPaid(ctx context.Context, ID string, Reply string) (bool, error) {
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
	cr, err := r.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}
	if cr.Status != "unchecked" {
		return false, errors.New("当前状态下不允许进行此操作")
	}

	SQL := `update cash_requests set status = $1, reply = $2 where id = $3`
	if _, err := r.db.Exec(SQL, "paid", Reply, ID); err != nil {
		return false, err
	}
	return true, nil
}

// 商家设置已完成
func (r *CashRequestRepository) SetFinished(ctx context.Context, ID string) (bool, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		return false, errors.New("只有商家能进行此操作")
	}
	cr, err := r.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}
	if cr.Status != "paid" {
		return false, errors.New("当前状态下不允许进行此操作")
	}
	SQL := `update cash_requests set status = $1 where id = $2`
	if _, err := r.db.Exec(SQL, "finished", ID); err != nil {
		return false, err
	}
	// 插入流水表记录
	logSQL := `INSERT INTO merchant_balance_logs (merchant_id,"inout","references",inout_type) VALUES ($1,$2,$3,$4);`
	if _, err = r.db.Exec(logSQL, gc.CurrentUser(ctx), -cr.Amount, ID, "withdrawal"); err != nil {
		return false, err
	}

	return true, nil
}
