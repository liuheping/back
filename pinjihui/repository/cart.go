package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    gc "pinjihui.com/pinjihui/context"
    "pinjihui.com/pinjihui/model"
    "github.com/rs/xid"
    "pinjihui.com/pinjihui/util"
    "database/sql"
    "errors"
    "fmt"
    "pinjihui.com/pinjihui/loader"
)

type CartRepository struct {
    BaseRepository
}

func NewCartRepository(db *sqlx.DB, log *logging.Logger) *CartRepository {
    return &CartRepository{BaseRepository{db: db, log: log}}
}

func (c *CartRepository) FindAll(ctx context.Context) ([]*model.CartItem, error) {
    return c.FindBy(ctx, nil)
}

func (c *CartRepository) TotalCount(ctx context.Context) (int, error) {
    gc.CheckAuth(ctx)
    return c.GetIntFormDB(`SELECT COALESCE(SUM(product_count), 0) FROM carts,rel_merchants_products r, products p WHERE p.id=carts.product_id AND carts.product_id=r.product_id AND carts.merchant_id=r.merchant_id AND user_id=$1`, gc.CurrentUser(ctx))
}

func (c *CartRepository) Save(ctx context.Context, newItem *model.CartItem) (*model.CartItem, error) {
    gc.CheckAuth(ctx)
    //检查是否存在商品
    _, err := loader.LoadProduct(ctx, newItem.ProductID, newItem.MerchantID)
    if err != nil {
        return nil, err
    }
    //限制加盟商购买加盟商的商品
    merchant, err := loader.LoadMerchant(ctx, newItem.MerchantID)
    if err != nil {
        return nil, err
    }
    if merchant.Type == model.ALLY && gc.IsAlly(ctx) {
        return nil, errors.New("ally can't buy product from ally")
    }
    newItem.UserID = *gc.CurrentUser(ctx)
    //check if has the same product in cart
    var oldItem model.CartItem
    querytSql := `SELECT id, product_count FROM carts WHERE user_id=$1 AND product_id=$2 AND merchant_id=$3`
    err = c.db.Get(&oldItem, querytSql, newItem.UserID, newItem.ProductID, newItem.MerchantID)
    if err != nil && err != sql.ErrNoRows {
        return nil, err
    }
    var sqlStr string
    if err == sql.ErrNoRows {
        //create. validate product id and merchant id
        newItem.ID = xid.New().String()
        sqlStr = util.NewSQLBuilder(newItem).Table("carts").InsertSQLBuild([]string{"Price"})
    } else {
        //update
        newItem.ProductCount += oldItem.ProductCount
        newItem.ID = oldItem.ID
        sqlStr = `UPDATE carts SET product_count=:product_count WHERE id=:id`
    }

    //check stock
    if err = c.checkStock(newItem, newItem.ProductCount); err != nil {
        return nil, err
    }
    //检查是否可以秒杀
    tx := c.db.MustBegin()
    _, err = L("spike").(*SpikeRepository).CanSpike(tx, newItem.ProductID, newItem.MerchantID, int(newItem.ProductCount))
    if err != nil {
        return nil, err
    }
    if _, err := tx.NamedExec(sqlStr, newItem); err != nil {
        return nil, err
    }
    tx.Commit()
    return newItem, nil
}

func (c *CartRepository) UpdateCount(ctx context.Context, id string, count int32) (*model.CartItem, error) {
    gc.CheckAuth(ctx)
    query := `SELECT id, product_id, product_count, merchant_id FROM carts WHERE id=$1`
    var oldItem model.CartItem
    err := c.db.Get(&oldItem, query, id)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("Invalid id: %s", id)
    }
    if err != nil {
        return nil, err
    }
    if err = c.checkStock(&oldItem, count); err != nil {
        return nil, err
    }
    sqlStr := `UPDATE carts SET product_count=$1 WHERE id=$2`

    if _, err := c.db.Exec(sqlStr, count, id); err != nil {
        return nil, err
    }
    oldItem.ProductCount = count
    return &oldItem, nil
}

func (c *CartRepository) checkStock(item *model.CartItem, newCount int32) error {
    stock, err := c.getStock(item.ProductID, item.MerchantID)
    if err == sql.ErrNoRows {
        return errors.New("product not exist")
    }
    if err != nil {
        return err
    }
    if stock == 0 {
        return errors.New("No Stock")
    }
    if stock < newCount {
        return errors.New("add too much products into cart")
    }
    return nil
}

func (c *CartRepository) getStock(productID, merchantID string) (int32, error) {
    stockSql := `SELECT stock FROM rel_merchants_products r,products WHERE 
r.product_id=products.id and products.deleted=false and products.is_sale=true AND 
product_id=$1 AND r.merchant_id=$2`
    var stock int32 = 0
    err := c.db.Get(&stock, stockSql, productID, merchantID)
    return stock, err
}

func (c *CartRepository) DeleteItem(ctx context.Context, id string) (bool, error) {
    gc.CheckAuth(ctx)
    delSql := `DELETE FROM carts WHERE id=$1 AND user_id=$2`
    r, err := c.db.Exec(delSql, id, *gc.CurrentUser(ctx))
    if err != nil {
        return false, err
    }
    if af, _ := r.RowsAffected(); af != 1 {
        return false, fmt.Errorf("item delete failed, id: %s", id)
    }
    return true, nil
}

func (c *CartRepository) FindByIDs(ctx context.Context, ids *[]string) ([]*model.CartItem, error) {
    return c.FindBy(ctx, ids)
}

func (c *CartRepository) FindBy(ctx context.Context, ids *[]string) ([]*model.CartItem, error) {
    gc.CheckAuth(ctx)
    items := make([]*model.CartItem, 0)
    var err error
    //加盟商看到的价格为二手价
    priceColumn := util.If(gc.IsAlly(ctx), "second_price", "retail_price").(string)
    sqlStr := fmt.Sprintf(`SELECT carts.id, carts.product_id, carts.product_count, carts.merchant_id, %s price FROM carts, rel_merchants_products r, products p WHERE p.id=carts.product_id AND carts.product_id=r.product_id AND carts.merchant_id=r.merchant_id AND user_id='%s'`, priceColumn, *gc.CurrentUser(ctx))
    if ids != nil {
        sqlStr += ` AND carts.id in (?)`
        query, args, err := sqlx.In(sqlStr, *ids)
        if err == nil {
            query = c.db.Rebind(query)
            //注意,即使返回空结果,err也为nil
            err = c.db.Select(&items, query, args...)
        }
    } else {
        err = c.db.Select(&items, sqlStr)
    }
    return items, err
}
