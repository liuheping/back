package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
	//gc "pinjihui.com/backend.pinjihui/context"
)

type PaymentRepository struct {
	BaseRepository
}

func NewPaymentRepository(db *sqlx.DB, log *logging.Logger) *PaymentRepository {
	return &PaymentRepository{BaseRepository{db: db, log: log}}
}

func (r *PaymentRepository) SavePayment(ctx context.Context, payment *model.Payment) (*model.Payment, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return nil, errors.New("只有管理员能增加修改支付方式")
	}
	if payment.ID == "" {
		payment.ID = xid.New().String()
		SQL = util.InsertSQLBuild(payment, "payments", []string{"Deleted", "Enabled", "Is_online", "Is_cod"})
	} else {
		if _, err := r.FindByID(ctx, payment.ID); err != nil {
			return nil, err
		}
		SQL = util.UpdateSQLBuild(payment, "payments", []string{"Deleted", "Enabled", "Is_online", "Is_cod"})
	}
	if _, err := r.db.NamedExec(SQL, payment); err != nil {
		return nil, err
	}
	pay, err := r.FindByID(ctx, payment.ID)
	if err != nil {
		return nil, err
	}
	return pay, nil
}

func (r *PaymentRepository) FindByID(ctx context.Context, ID string) (*model.Payment, error) {
	payment := &model.Payment{}
	SQL := `SELECT * FROM payments WHERE id = $1`
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepository) FindAll(ctx context.Context) ([]*model.Payment, error) {
	payments := make([]*model.Payment, 0)
	SQL := `SELECT * FROM payments where deleted=false ORDER BY sort_order ASC ;`
	if err := r.db.Select(&payments, SQL); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return false, errors.New("只有管理员能删除支付方式")
	}
	SQL := `update payments set deleted=not deleted where id=$1`
	if _, err := r.db.Exec(SQL, ID); err != nil {
		return false, err
	}
	return true, nil
}
