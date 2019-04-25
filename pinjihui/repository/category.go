package repository

import (
    "fmt"

    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/util"
)

type CategoryRepository struct {
    BaseRepository
}

const CATE_TABLE = "categories"

func NewCategoryRepository(db *sqlx.DB, log *logging.Logger) *CategoryRepository {
    return &CategoryRepository{BaseRepository{db: db, log: log}}
}

func (r *CategoryRepository) List(parent_id *string) ([]*model.Category, error) {
    whereP := util.WhereNullable(parent_id)
    sqlStr := fmt.Sprintf(`SELECT id, parent_id, name, thumbnail FROM %s WHERE deleted=false AND enabled=true AND parent_id %s ORDER BY sort_order`,
        CATE_TABLE, whereP)
    fmt.Println(sqlStr)
    ms := make([]*model.Category, 0)
    if err := r.db.Select(&ms, sqlStr); err != nil {
        return nil, err
    }
    return ms, nil
}

func (r *CategoryRepository) GetTree(machineTypeSeries *string, merchantId *string, parentId *string) (*model.Category, error) {
    var sqlStr string
    ms := make([]*model.Category, 0)
    args := make([]interface{}, 0)
    if merchantId == nil {
        if machineTypeSeries == nil {
            if parentId != nil {
                args = append(args, parentId)
                sqlStr = `SELECT id, parent_id, name, thumbnail,is_common FROM categories WHERE parent_id=$1 AND enabled = true AND id IN (
                        SELECT DISTINCT category_id
                        FROM products
                        WHERE deleted = false AND is_sale = true)`
            } else {
                sqlStr = `WITH RECURSIVE cateTree AS (
                      SELECT
                        id,
                        parent_id,
                        name,
                        thumbnail,
                        is_common,
                        sort_order
                      FROM categories
                      WHERE enabled = true AND id IN (
                        SELECT DISTINCT category_id
                        FROM products
                        WHERE deleted = false AND is_sale = true)
                      UNION ALL SELECT
                                  DISTINCT
                                  c2.id,
                                  c2.parent_id,
                                  c2.name,
                                  c2.thumbnail,
                                  c2.is_common,
                                  c2.sort_order
                                FROM categories c2, cateTree t
                                WHERE c2.id = t.parent_id)
                    SELECT DISTINCT *
                    FROM cateTree
                    ORDER BY sort_order ASC`
            }
        } else {
            args = append(args, machineTypeSeries)
            sqlStr = `
                WITH RECURSIVE cateTree AS (
                  SELECT DISTINCT
                    c1.id,
                    c1.parent_id,
                    c1.name,
                    c1.is_common,
                    c1.thumbnail
                  FROM categories c1
                  WHERE c1.id IN (SELECT DISTINCT p.category_id FROM products p JOIN brand_series bs ON p.machine_types && bs.machine_types WHERE bs.id=$1 AND is_sale=true AND deleted=false) AND c1.enabled=true
                  UNION ALL SELECT DISTINCT
                              c2.id,
                              c2.parent_id,
                              c2.name,
                              c2.is_common,
                              c2.thumbnail
                            FROM categories c2, cateTree t
                            WHERE c2.id = t.parent_id
                )
                SELECT DISTINCT * FROM cateTree;`
        }
    } else {
        args = append(args, merchantId)
        sqlStr = `
                    WITH RECURSIVE cateTree AS (
                    SELECT id, parent_id, name, thumbnail FROM categories WHERE id in (SELECT DISTINCT category_id FROM products WHERE id in (SELECT product_id FROM rel_merchants_products WHERE merchant_id = $1)) AND enabled=true 
                      UNION ALL SELECT DISTINCT
                                  c2.id,
                                  c2.parent_id,
                                  c2.name,
                                  c2.thumbnail
                                FROM categories c2, cateTree t
                                WHERE c2.id = t.parent_id
                    )
                    SELECT DISTINCT * FROM cateTree`
    }
    if err := r.db.Select(&ms, sqlStr, args...); err != nil {
        return nil, err
    }
    return buildTree(&ms), nil
}

func buildTree(categories *[]*model.Category) (root *model.Category) {
    tempMap := make(map[string]*model.Category)
    for _, v := range *categories {
        tempMap[v.ID] = v
    }
    root = &model.Category{Children: make([]*model.Category, 0)}
    for _, v := range *categories {
        if v.ParentId == nil {
            root.Children = append(root.Children, v)
        } else if _, ok := tempMap[*v.ParentId]; !ok {
            root.Children = append(root.Children, v)
        } else {
            tempMap[*v.ParentId].Children = append(tempMap[*v.ParentId].Children, v)
        }
    }
    return root
}
