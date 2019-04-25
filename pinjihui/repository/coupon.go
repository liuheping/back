package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
    gc "pinjihui.com/pinjihui/context"
    "golang.org/x/net/context"
    "fmt"
    "pinjihui.com/pinjihui/util"
)

type CouponRepository struct {
    BaseRepository
}

func NewCouponRepository(db *sqlx.DB, log *logging.Logger) *CouponRepository {
    return &CouponRepository{BaseRepository{db: db, log: log}}
}

func (c *CouponRepository) List(ctx context.Context, condition string, first int, after *string) ([]*model.Coupon, error) {
    gc.CheckAuth(ctx)
    coupons := make([]*model.Coupon, 0)
    builder := util.NewSQLBuilder(&coupons).Table("user_coupons").WhereRowWithHolderPlace("user_id=$1", gc.CurrentUser(ctx)).WhereRow(condition).
        OrderBy("created_at", "DESC").Limit(first, nil)
    if after != nil {
        builder.WhereRowWithHolderPlace(`created_at < (SELECT created_at FROM user_coupons WHERE id=$2)`, after)
    }
    fmt.Println(builder.BuildQuery())
    if err := c.db.Select(&coupons, builder.BuildQuery(), builder.Args...); err != nil {
        return nil, err
    }
    return coupons, nil
}

func (c *CouponRepository) FindAvailable(ctx context.Context, first int, after *string) ([]*model.Coupon, error) {
    return c.List(ctx, c.GetConditionBy("available"), first, after)
}

func (c *CouponRepository) FindByStatus(ctx context.Context, status string, first int, after *string) ([]*model.Coupon, error) {
    return c.List(ctx, c.GetConditionBy(status), first, after)
}

func (c *CouponRepository) GetConditionBy(status string) string {
    switch status {
    case "available":
        return "used=false AND expired_at >= CURRENT_DATE AND start_at<=CURRENT_DATE"
    case "used":
        return "used=true"
    case "expired":
        return "used=false AND expired_at<CURRENT_DATE"
    case "not_used":
        return "used=false AND expired_at>=CURRENT_DATE"
    default:
        return ""
    }
}

func (c *CouponRepository) Count(ctx context.Context, status string) (int, error) {
    sqls := fmt.Sprintf(`SELECT COUNT(*) FROM user_coupons WHERE user_id='%s' AND %s`, *gc.CurrentUser(ctx), c.GetConditionBy(status))
    return c.GetIntFormDB(sqls)
}

func (c *CouponRepository) FindAvailableByIDs(ctx context.Context, ids []string) ([]*model.Coupon, error) {
    gc.CheckAuth(ctx)
    coupons := make([]*model.Coupon, 0)

    couponSQL := fmt.Sprintf(`SELECT id, value, limit_amount FROM user_coupons WHERE id in (?) AND user_id='%s' AND used=false AND expired_at >= CURRENT_DATE AND start_at<=CURRENT_DATE AND merchant_id IS NUll`, *gc.CurrentUser(ctx))
    query, args, err := sqlx.In(couponSQL, ids)
    if err != nil {
        return nil, err
    }
    query = c.db.Rebind(query)
    fmt.Println(query, args)
    err = c.db.Select(&coupons, query, args...)
    if err != nil {
        c.log.Errorf("Error in retrieving coupons : %v", err)
        return nil, err
    }
    return coupons, nil
}
