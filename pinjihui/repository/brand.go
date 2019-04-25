package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
    "fmt"
    "pinjihui.com/pinjihui/util"
)

type BrandRepository struct {
    BaseRepository
}

func NewBrandRepository(db *sqlx.DB, log *logging.Logger) *BrandRepository {
    return &BrandRepository{BaseRepository{db: db, log: log}}
}

func (r *BrandRepository) List(brandType string, first *int32, merchantId *string, cateId *string) ([]*model.Brand, error) {
    var SQL string
    ms := make([]*model.BrandDB, 0)
    args := []interface{}{first}
    if merchantId == nil {
        if brandType == "excavator" {
            whereCate := ""
            if cateId != nil {
                whereCate = ` AND category_id IN (WITH RECURSIVE cateTree AS (
  SELECT id FROM categories WHERE id =$2 UNION ALL SELECT c.id FROM categories c,cateTree WHERE c.parent_id=cateTree.id
)
SELECT id from cateTree) `
                args = append(args, cateId)
            }
            SQL = fmt.Sprintf(`SELECT
            id, name, thumbnail, description,array_agg(mt) machine_types
            FROM (select
                b.*,
                unnest(b.machine_types) mt
                from brands b
                where type = 'excavator' and enabled=true and deleted=false ) a
                where a.mt in (select DISTINCT unnest(machine_types)
                from products
                where deleted = false AND is_sale=true %s) 
                group by id, name, thumbnail, description,sort_order order by sort_order limit $1`, whereCate)
        } else {
            SQL = `SELECT
                  id,
                  name,
                  thumbnail,
                  description,
                  machine_types
                FROM brands
                WHERE type = 'part' AND deleted = false AND enabled = true AND id IN (SELECT distinct brand_id
                                                                                      from products
                                                                                      where brands.deleted = false and is_sale = true)
                order by sort_order limit $1`
        }
        if err := r.db.Select(&ms, SQL, args...); err != nil {
            return nil, err
        }
    } else {
        SQL = `SELECT id, name, thumbnail, description,machine_types FROM brands WHERE id in (SELECT brand_id FROM products WHERE id in (SELECT product_id FROM rel_merchants_products WHERE merchant_id = $1)) AND type=$2 AND deleted=false AND enabled=true order by sort_order limit $3 `
        if err := r.db.Select(&ms, SQL, merchantId, brandType, first); err != nil {
            return nil, err
        }
    }

    brands := model.GetBrands(&ms)
    return brands, nil
}

func (b *BrandRepository) FindSeries(brandId string) ([]*model.BrandSeries, error) {
    series := make([]*model.BrandSeries, 0)
    err := b.db.Select(&series, `SELECT id, series, image, machine_types FROM brand_series WHERE brand_id=$1`, brandId)
    return series, err
}

func (b *BrandRepository) FindMachineTypes(seriesId string) ([]string, error) {
    var types string
    err := b.db.Get(&types, `SELECT machine_types FROM brand_series WHERE id=$1`, seriesId)
    return util.ParseArray(&types), err
}
