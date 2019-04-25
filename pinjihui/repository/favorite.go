package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
    "golang.org/x/net/context"
    gc "pinjihui.com/pinjihui/context"
    "pinjihui.com/pinjihui/loader"
    "github.com/rs/xid"
    "pinjihui.com/pinjihui/util"
    "errors"
    "fmt"
    "database/sql"
)

type FavoriteRepository struct {
    BaseRepository
}

func NewFavoriteRepository(db *sqlx.DB, log *logging.Logger) *FavoriteRepository {
    return &FavoriteRepository{BaseRepository{db: db, log: log}}
}

func (f *FavoriteRepository) Add(ctx context.Context, merchantID string, productID *string) (*model.Favorite, error) {
    gc.CheckAuth(ctx)
    var err error
    if productID == nil {
        //收藏店铺,检查店铺是否存在
        _, err = loader.LoadMerchant(ctx, merchantID)
    } else {
        _, err = loader.LoadProduct(ctx, *productID, merchantID)
    }
    if err != nil {
        return nil, err
    }
    hasOne, err := f.HasOne(ctx, merchantID, productID)
    if err != nil {
        return nil, err
    }
    if hasOne {
        return nil, errors.New("favorite already exist")
    }
    favorite := &model.Favorite{
        ID:         xid.New().String(),
        UserID:     *gc.CurrentUser(ctx),
        MerchantID: merchantID,
        ProductID:  productID,
    }
    sqls := util.NewSQLBuilder(favorite).InsertSQLBuild(nil)
    _, err = f.db.NamedExec(sqls, favorite)
    if err != nil {
        return nil, err
    }
    return favorite, nil
}

func (f *FavoriteRepository) HasOne(ctx context.Context, merchantID string, productID *string) (bool, error) {
    args, where := f.GenerateWhereByMerchantProductID(ctx, merchantID, productID)
    sqls := `SELECT count(*) FROM favorites WHERE ` + where
    var count int
    err := f.db.Get(&count, sqls, args...)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (f *FavoriteRepository) GenerateWhereByMerchantProductID(ctx context.Context, merchantID string, productID *string) ([]interface{}, string) {
    var where = "user_id=$1 AND merchant_id=$2 AND product_id "
    args := make([]interface{}, 0)
    args = append(args, *gc.CurrentUser(ctx))
    args = append(args, merchantID)
    if productID == nil {
        where += "IS NULL"
    } else {
        args = append(args, *productID)
        where += "=$3"
    }
    return args, where
}

func (f *FavoriteRepository) FindByPMID(ctx context.Context, merchantID string, productID *string) (*model.Favorite, error) {
    fav := model.Favorite{}
    args, where := f.GenerateWhereByMerchantProductID(ctx, merchantID, productID)
    sqls := util.NewSQLBuilder(&fav).WhereRow(where).BuildQuery()
    err := f.db.Get(&fav, sqls, args...)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    return &fav, err
}

const productType = "product"

func (f *FavoriteRepository) List(ctx context.Context, first int32, after *string, favoriteType string) ([]*model.Favorite, error) {
    gc.CheckAuth(ctx)
    list := make([]*model.Favorite, 0)
    builder := util.NewSQLBuilder(list)
    f.BuildWhere(builder, favoriteType)
    args := []interface{}{gc.CurrentUser(ctx)}
    if after != nil {
        builder.WhereRow("created_at<(SELECT created_at FROM favorites WHERE id = $2)")
        args = append(args, after)
    }
    sqls := builder.OrderBy("created_at", "DESC").Limit(int(first), nil).BuildQuery()
    fmt.Println(sqls)
    err := f.db.Select(&list, sqls, args...)
    if err != nil {
        return nil, err
    }
    return list, nil
}

func (f *FavoriteRepository) Count(ctx context.Context, favoriteType string) (int, error) {
    gc.CheckAuth(ctx)
    var count int
    builder := util.NewSQLBuilder(nil).Select("count(*)").Table("favorites")
    f.BuildWhere(builder, favoriteType)
    sqls := builder.BuildQuery()
    err := f.db.Get(&count, sqls, *gc.CurrentUser(ctx))
    return count, err
}

func (f *FavoriteRepository) BuildWhere(builder *util.SQLBuilder, favoriteType string) {
    where := util.If(favoriteType == productType, "product_id IS NOT NULL", "product_id IS NULL").(string)
    builder.WhereRow("user_id=$1").WhereRow(where)
}

func (f *FavoriteRepository) Remove(ctx context.Context, id string) error {
    gc.CheckAuth(ctx)
    sqls := `DELETE FROM favorites WHERE id=$1 AND user_id=$2`
    result, err := f.db.Exec(sqls, id, *gc.CurrentUser(ctx))
    if err != nil {
        return err
    }
    af, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if af == 0 {
        return gc.ErrNoRecord
    }
    return nil
}
