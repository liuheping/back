package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
    "errors"
    "pinjihui.com/pinjihui/util"
    "fmt"
    "golang.org/x/net/context"
    "github.com/rs/xid"
    gc "pinjihui.com/pinjihui/context"
    "database/sql"
)

type MerchantRepository struct {
    BaseRepository
}

func NewMerchantRepository(db *sqlx.DB, log *logging.Logger) *MerchantRepository {
    return &MerchantRepository{BaseRepository{db: db, log: log}}
}

func (r *MerchantRepository) List(first int, offset int, orderBy string, position *model.Location) ([]*model.Merchant, error) {
    merchantDBS := make([]*model.MerchantDB, 0)
    builder := util.NewSQLBuilder(merchantDBS).Table("merchant_profiles").Alias("m").
        Join("users u", "u.id=m.user_id").
        WhereRow("u.type='ally'")
    var sort string
    args := make([]interface{}, 0)
    if orderBy == "distance" {
        if position == nil {
            return nil, errors.New("need position if order by distance")
        }
        builder.AddSelect("earth_distance(ll_to_earth(m.lat, m.lng), ll_to_earth($1, $2)) distance")
        args = append(args, position.Lat, position.Lng)
        sort = "ASC"
    } else if orderBy == "sales_volume" {
        builder.LeftJoin("rel_merchants_products r", "r.merchant_id=m.user_id").
            AddSelect("COALESCE(SUM(r.sales_volume), 0) sales_volume").
            GroupBy("m.user_id")
        sort = "DESC"
    } else if orderBy == "comment_quantity" {
        builder.LeftJoin("comments c", "c.merchant_id=m.user_id").
            AddSelect("count(c.*) comment_quantity").
            GroupBy("m.user_id")
        sort = "DESC"
    } else {
        //can't be here
        panic("?")
    }
    query := builder.AddSelect("'ally' AS type").OrderBy(orderBy, sort).Limit(first, &offset).BuildQuery()
    fmt.Println(query)
    if err := r.db.Select(&merchantDBS, query, args...); err != nil {
        return nil, err
    }
    merchants := make([]*model.Merchant, len(merchantDBS))
    for i, v := range merchantDBS {
        if err := v.ParseAddrs(); err != nil {
            return nil, err
        }
        merchants[i] = &v.Merchant
    }
    return merchants, nil
}

func (r *MerchantRepository) Count() (int, error) {
    return r.GetIntFormDB(`SELECT COUNT(*) FROM merchant_profiles m, users u WHERE m.user_id=u.id AND u.type='ally'`)
}

func (r *MerchantRepository) GetDistance(id string, l *model.Location) (*float64, error) {
    var d *float64
    sqls := `SELECT earth_distance(ll_to_earth(m.lat, m.lng), ll_to_earth($1, $2)) distance FROM merchant_profiles m WHERE user_id=$3`
    err := r.db.Get(&d, sqls, l.Lat, l.Lng, id)
    return d, err
}

func (r *MerchantRepository) GetSalesVolume(id string) (int, error) {
    sqls := `SELECT COALESCE(SUM(op.product_number), 0) FROM order_products op JOIN orders o ON op.order_id=o.id AND o.merchant_id=$1 AND o.status='paid'`
    return r.GetIntFormDB(sqls, id)
}

func (r *MerchantRepository) GetCollectedQuantity(id string) (int, error) {
    sqls := `SELECT count(*) FROM favorites WHERE merchant_id=$1 AND product_id IS NULL;`
    return r.GetIntFormDB(sqls, id)
}

func (r *MerchantRepository) GetCommentQuantity(id string) (int, error) {
    sqls := `SELECT count(*) FROM comments WHERE merchant_id=$1;`
    return r.GetIntFormDB(sqls, id)
}

func (r *MerchantRepository) FindByID(id string) (*model.Merchant, error) {
    m := &model.MerchantDB{}
    sqlStr := util.NewSQLBuilder(m).Table("merchant_profiles").Alias("m").Join("users u", "u.id=m.user_id").
        AddSelect("u.type").WhereRow("user_id = $1").
        BuildQuery()
    err := r.db.Get(m, sqlStr, id)
    if err != nil {
        m.ParseAddrs()
    }
    return &m.Merchant, err
}

func (r *MerchantRepository) FindWaitersByMerchantId(id string) ([]string, error) {
    ids := make([]string, 0)
    query := `SELECT waiter_id FROM waiters WHERE merchant_id=$1 AND waiter_id IS NOT NULL`
    err := r.db.Select(&ids, query, id)
    return ids, err
}

func (r *MerchantRepository) Purchase(ctx context.Context, productID string) error {
    gc.CheckAuth(ctx)
    if gc.User(ctx).Type != model.UTAgent {
        return errors.New("你没有权限进货")
    }
    id := xid.New().String()
    query := `INSERT INTO rel_agents_products(id, agent_id, product_id) VALUES ($1, $2, $3) ON CONFLICT(agent_id, product_id) DO UPDATE SET is_sale=TRUE`
    _, err := r.db.Exec(query, id, gc.CurrentUser(ctx), productID)
    return err
}

func (r *MerchantRepository) UnPurchase(ctx context.Context, productID string) error {
    gc.CheckAuth(ctx)
    if gc.User(ctx).Type != model.UTAgent {
        return errors.New("你没有权限下架")
    }
    query := `UPDATE rel_agents_products SET is_sale=FALSE WHERE agent_id=$1 AND product_id=$2`
    _, err := r.db.Exec(query, gc.CurrentUser(ctx), productID)
    return err
}

func (r *MerchantRepository) IsPurchased(ctx context.Context, productID string) (bool, error) {
    gc.CheckAuth(ctx)
    var isPurchased bool
    query := `SELECT is_sale FROM rel_agents_products WHERE agent_id=$1 AND product_id=$2`
    err := r.db.Get(&isPurchased, query, gc.CurrentUser(ctx), productID)
    if err == sql.ErrNoRows {
        err = nil
    }
    return isPurchased, err
}