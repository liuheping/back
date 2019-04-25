package loader

import (
    "pinjihui.com/pinjihui/model"
    "fmt"
    "pinjihui.com/pinjihui/util"
    "golang.org/x/net/context"
    "github.com/jmoiron/sqlx"
    "database/sql"
    "github.com/op/go-logging"
)

func FindProductWithStockByIDs(ctx context.Context, pms []*model.PM) ([]*model.PaMCPair, error) {
    products := make([]*model.PaMCPair, 0)
    where := "("
    for i, v := range pms {
        if i != 0 {
            where += " OR "
        }
        where += fmt.Sprintf(" r.merchant_id='%s' AND r.product_id='%s' ", v.MerchantID, v.ProductID)
    }
    where += ")"
    sqlStr := util.NewSQLBuilder(model.Product{}).AddSelect("r.merchant_id", "r.stock", "r.retail_price", "r.origin_price", "r.sales_volume", "r.is_sale", "r.view_volume").
        Join("rel_merchants_products r", "products.id=r.product_id").
        WhereRow("deleted=false").WhereRow(where).
        BuildQuery()

    if err := ctx.Value("db").(*sqlx.DB).Select(&products, sqlStr); err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("err when get product with err: %v, sql: %s", err, sqlStr)
        return nil, err
    }

    return products, nil
}

func FindMerchantsByIds(ctx context.Context, ids []string)  ([]*model.Merchant, error) {
    merchantDBS := make([]*model.MerchantDB, 0)
    sqlStr := util.NewSQLBuilder(merchantDBS).Table("merchant_profiles").Alias("m").Join("users u", "u.id=m.user_id").
        AddSelect("u.type").WhereRow("user_id in (?)").
        BuildQuery()
    query, args, err := sqlx.In(sqlStr, ids)
    if err != nil {
        return nil, err
    }
    db := ctx.Value("db").(*sqlx.DB)
    query = db.Rebind(query)
    if err = db.Select(&merchantDBS, query, args...); err != nil {
        return nil, err
    }
    merchants := make([]*model.Merchant, len(merchantDBS))
    for i, v := range merchantDBS {
        if err = v.ParseAddrs(); err != nil {
            return nil, err
        }
        merchants[i] = &v.Merchant
    }
    return merchants, nil
}

func FindBrandsByIds(ctx context.Context, ids []string) ([]*model.Brand, error) {
    sqlStr := `SELECT id, name, thumbnail, description, machine_types FROM brands WHERE id IN (?)`
    query, args, err := sqlx.In(sqlStr, ids)
    if err != nil {
        return nil, err
    }
    db := ctx.Value("db").(*sqlx.DB)
    query = db.Rebind(query)
    dbBrands := make([]*model.BrandDB, 0)
    if err = db.Select(&dbBrands, query, args...); err != nil && err != sql.ErrNoRows {
        return nil, err
    }
    return model.GetBrands(&dbBrands), nil
}

func FindUserByID(ctx context.Context,ID string) (*model.User, error) {
    user := &model.User{}

    userSQL := `SELECT * FROM users WHERE id = $1`
    db := ctx.Value("db").(*sqlx.DB)
    row := db.QueryRowx(userSQL, ID)
    err := row.StructScan(user)
    if err != nil {
        return nil, err
    }
    return user, nil
}
func FindTheClosestMerchant(ctx context.Context, productIds []string, location *model.Location) ([]*model.MerchantWithStock, error) {
    merchants := make([]*model.MerchantWithStock, 0)

    merchantsSQL := fmt.Sprintf(`select
        user_id,
        company_name,
        company_address,
        delivery_address,
        product_id,
        stock,
        retail_price,
  distance
from (SELECT *,
        earth_distance(ll_to_earth(m.lat, m.lng), ll_to_earth(%f, %f)) distance,
        rank()
        over (
          partition by product_id
          order by earth_distance(ll_to_earth(m.lat, m.lng), ll_to_earth(%f, %f)) )
      FROM merchant_profiles m, rel_merchants_products mp,users u
      where m.user_id = mp.merchant_id
            AND u.id=m.user_id AND u.type='ally' AND mp.product_id in (?)) b
where rank = 1`, location.Lat, location.Lng, location.Lat, location.Lng)
    query, args, err := sqlx.In(merchantsSQL, productIds)
    if err != nil {
        fmt.Errorf("sqlx.In(merchantsSQL, productIds) err: %v", err)
        return nil, err
    }
    db := ctx.Value("db").(*sqlx.DB)
    query = db.Rebind(query)
    err = db.Select(&merchants, query, args...)

    if err != nil {
        return nil, err
    }
    for _, v := range merchants {
        if err = v.ParseAddrs(); err != nil {
            return nil, err
        }
    }
    return merchants, nil

}