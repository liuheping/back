package repository

import (
    "database/sql"
    "pinjihui.com/pinjihui/model"
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    gc "pinjihui.com/pinjihui/context"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/util"
    "fmt"
    "golang.org/x/net/context"
    "strings"
)

type ProductRepository struct {
    BaseRepository
}

func NewProductRepository(db *sqlx.DB, log *logging.Logger) *ProductRepository {
    return &ProductRepository{BaseRepository{db: db, log: log}}
}

func (p *ProductRepository) FindByID(ID string) (*model.Product, error) {
    product := &model.Product{}

    productSQL := util.NewSQLBuilder(model.Product{}).WhereRow("id=$1 AND deleted=false").BuildQuery()
    row := p.db.QueryRowx(productSQL, ID)
    err := row.StructScan(product)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    if err != nil {
        p.log.Errorf("Error in retrieving product : %v", err)
        return nil, err
    }
    return product, nil
}

func (p *ProductRepository) FindProductSpecs(product *model.PaMCPair) ([]*model.ProductSpec, error) {
    items := make([]*model.ProductSpec, 0)
    sqls := `SELECT id, spec_1, spec_2, retail_price, stock FROM products, rel_merchants_products r WHERE products.id=r.product_id AND r.merchant_id=$1 AND parent_id = $2 AND products.deleted=FALSE AND r.is_sale=TRUE`
    err := p.db.Select(&items, sqls, product.MerchantID, util.GetString(product.ParentID, product.ID))
    if err != nil && err != sql.ErrNoRows {
        return nil, err
    }
    return items, nil
}

func (p *ProductRepository) FindByIDs(ids []string) ([]*model.Product, error) {

    products := make([]*model.Product, 0)
    sqlStr := util.NewSQLBuilder(products).WhereRow(`id in (?) AND delete=false`).BuildQuery()
    query, args, err := sqlx.In(sqlStr, ids)
    if err != nil {
        return nil, err
    }
    query = p.db.Rebind(query)
    if err = p.db.Select(&products, query, args...); err != nil {
        return nil, err
    }

    return products, nil
}

type ProductSearchInput struct {
    Key               *string     `db:"-"`
    Brand             *graphql.ID `db:"products.brand_id"`
    Category          *graphql.ID `db:"-"`
    Merchant          *graphql.ID `db:"r.merchant_id"`
    MachineType       *string     `db:"products.machine_types" func:"any"`
    MachineTypeSeries *graphql.ID `db:"-"`
}

type ProductSortInput struct {
    OrderBy  string
    Sort     *string
    Position *model.Location
}

func (p *ProductRepository) Search(ctx context.Context, fetchSize int, offset int, search *ProductSearchInput, sort *ProductSortInput) ([]*model.PaMCPair, error) {
    if sort != nil && sort.OrderBy == distance && sort.Position == nil {
        return nil, gc.ErrInvalidParam
    }
    products := make([]*model.PaMCPair, 0)
    builder := util.NewSQLBuilder(model.Product{}).
        Join("rel_merchants_products AS r", "products.id=r.product_id").
        AddSelect("r.merchant_id", "r.stock", "r.retail_price", "r.origin_price", "r.sales_volume", "r.is_sale", "r.view_volume").
        AddSelect(`row_number() over (partition by r.merchant_id,parent_id order by retail_price) AS rank`)
    p.searchWhere(search, builder)
    buildSortSQL(sort, builder)
    productSQL := util.NewSQLBuilder(products).Table(fmt.Sprintf("(%s) a", builder.BuildQuery())).
        WhereRow("rank=1 OR parent_id IS NULL").
        Limit(fetchSize, &offset).BuildQuery()
    fmt.Println(productSQL)
    err := p.db.Select(&products, productSQL, builder.Args...)
    if err != nil {
        return nil, err
    }
    return products, nil
}

const (
    synthesis   = "synthesis"
    SalesVolume = "sales_volume"
    price       = "price"
    distance    = "distance"
    createdAt   = "created_at"
)

func buildSortSQL(sort *ProductSortInput, builder *util.SQLBuilder) {
    if sort == nil || sort.OrderBy == synthesis {
        return
    }
    var sortDire string
    orderBy := util.If(sort.OrderBy == price, "retail_price", sort.OrderBy).(string)
    if sort.OrderBy == SalesVolume || sort.OrderBy == createdAt {
        orderBy = "r." + sort.OrderBy
        sortDire = "DESC"
    } else if sort.OrderBy == distance {
        builder.Join("merchant_profiles m", "r.merchant_id = m.user_id")
        orderBy = fmt.Sprintf("earth_distance(ll_to_earth(m.lat, m.lng), ll_to_earth(%f, %f)) ", sort.Position.Lat, sort.Position.Lng)
        sortDire = "ASC"
    } else {
        sortDire = util.GetString(sort.Sort, " ASC ")
    }
    builder.OrderBy(orderBy, sortDire)
}

func (p *ProductRepository) searchWhere(search *ProductSearchInput, builder *util.SQLBuilder) {
    builder.WhereStruct(search, true).WhereRow("deleted=false AND r.is_sale=true")
    if search != nil {
        if search.Key != nil {
            builder.WhereRowWithHolderPlace(fmt.Sprintf("to_tsvector('zhparser', name) @@ to_tsquery('zhparser', $%d)", len(builder.Args)+1), strings.Replace(*search.Key, " ", "", -1))
        }
        if search.MachineTypeSeries != nil {
            builder.Join("brand_series bs", "products.machine_types && bs.machine_types").
                WhereRowWithHolderPlace(fmt.Sprintf("bs.id=$%d", len(builder.Args)+1), search.MachineTypeSeries)
        }
        if search.Category != nil {
            builder.WhereRowWithHolderPlace(fmt.Sprintf(`products.category_id IN (WITH RECURSIVE cateTree AS (
            SELECT id FROM categories WHERE id =$%d UNION ALL SELECT c.id FROM categories c,cateTree WHERE c.parent_id=cateTree.id
            )
            SELECT id from cateTree)`, len(builder.Args)+1), search.Category)
        }
    }
}

func (p *ProductRepository) Count(ctx context.Context, c *ProductSearchInput) (int, error) {
    builder := util.NewSQLBuilder(nil).Table("products").
        Select(`parent_id, row_number() over (partition by r.merchant_id,parent_id order by retail_price) AS rank`).
        Join("rel_merchants_products AS r", "products.id=r.product_id")
    p.searchWhere(c, builder)
    productSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) a WHERE rank=1 OR parent_id IS NULL", builder.BuildQuery())
    fmt.Println(productSQL)
    return p.GetIntFormDB(productSQL, builder.Args...)
}

func (p *ProductRepository) FindImagesById(ID string) ([]*model.ProductImage, error) {
    images := make([]*model.ProductImage, 0)
    sqls := util.NewSQLBuilder(images).WhereRow("product_id = $1").
        OrderBy("created_at", "ASC").BuildQuery()
    if err := p.db.Select(&images, sqls, ID); err != nil {
        return nil, err
    }
    return images, nil
}

func (p *ProductRepository) FindRelMerchantsProducts(productID, merchantID string) (*model.RelMerchantProduct, error) {
    rel := model.RelMerchantProduct{}
    sqls := util.NewSQLBuilder(&rel).Table("rel_merchants_products").
        WhereRow(`product_id=$1 AND merchant_id=$2`).
        BuildQuery()
    row := p.db.QueryRowx(sqls, productID, merchantID)
    err := row.StructScan(&rel)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    if err != nil {
        return nil, err
    }
    return &rel, nil
}

func (p *ProductRepository) GetRelMerchantsProducts(product *model.Product, merchantID string) (*model.PaMCPair, error) {
    rel, err := p.FindRelMerchantsProducts(product.ID, merchantID)
    if err != nil {
        return nil, err
    }
    return &model.PaMCPair{*product, *rel}, nil
}

//根据用户角色不同得到不同的价格
func (p *ProductRepository) GetPrice(product *model.PaMCPair, ctx context.Context) float64 {
    if ctx.Value("is_authorized").(bool) && gc.User(ctx).Type == model.ALLY {
        return product.SecondPrice
    }
    return product.Price
}

func (p *ProductRepository) AddViewCount(productID, merchantID string) error {
    query := `UPDATE rel_merchants_products SET view_volume=view_volume+1 WHERE product_id=$1 AND merchant_id=$2`
    _, err := p.db.Exec(query, productID, merchantID)
    return err
}