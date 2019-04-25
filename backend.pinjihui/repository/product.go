package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type ProductRepository struct {
	BaseRepository
}

func NewProductRepository(db *sqlx.DB, log *logging.Logger) *ProductRepository {
	return &ProductRepository{BaseRepository{db: db, log: log}}
}

func (r *ProductRepository) GetPriceRatio(ctx context.Context, code string) (*float64, error) {
	con, err := L("config").(*ConfigRepository).FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	ratio, err := strconv.ParseFloat(con.Value, 64)
	if err != nil {
		return nil, err
	}
	return &ratio, nil
}

// 上传商品（只有供货商能上传）
func (r *ProductRepository) Save(ctx context.Context, pro *model.ProductARR) (*model.Product, error) {
	// 查找品牌,获取溢价比率
	brand, err := L("brand").(*BrandRepository).FindByID(ctx, pro.Product.Brand_id)
	if err != nil {
		return nil, err
	}
	secondpriceratio := &brand.Second_price_ratio
	retailpriceratio := &brand.Retail_price_ratio
	// 获取用户类型
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "provider" {
		return nil, errors.New("只有供货商能上传商品")
	}
	pro.Product.Merchant_id = *gc.CurrentUser(ctx)
	pid := xid.New().String()

	if pro.Spec == nil {
		if pro.BatchPrice == nil || pro.Stock == nil {
			return nil, errors.New("价格或者库存不能为空")
		}
		pro.Product.ID = pid
		pro.Product.Type = "simple"
		pro.Product.Batch_price = pro.BatchPrice
		x := *pro.BatchPrice * (float64(1) + *secondpriceratio)
		pro.Product.Second_price = &x
		// 插入商品表
		SQL := util.InsertSQLBuild(pro.Product, "products", []string{"Is_sale", "Created_at", "Updated_at", "Deleted", "Recommended", "Spec_1_name", "Spec_2_name"})
		if _, err := r.db.NamedExec(SQL, pro.Product); err != nil {
			return nil, err
		}
		// 插入商家商品关联表
		RelSQL := `INSERT INTO rel_merchants_products (product_id,merchant_id,stock,retail_price) values ($1,$2,$3,$4)`
		RetailPrice := *pro.BatchPrice * (float64(1) + *retailpriceratio)
		_, err = r.db.Exec(RelSQL, pro.Product.ID, pro.Product.Merchant_id, pro.Stock, RetailPrice)
		if err != nil {
			return nil, err
		}
		// 插入图片
		ImageSQL := `INSERT INTO product_images (id, product_id, big_image) values ($1,$2,$3)`
		if pro.Images != nil {
			for _, image := range *pro.Images {
				_, err = r.db.Exec(ImageSQL, xid.New().String(), pro.Product.ID, image)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		//找最小价格
		min := (*pro.Spec)[0].Batch_price
		for _, v := range *pro.Spec {
			if min > v.Batch_price {
				min = v.Batch_price
			}
		}
		// 先插入父商品
		pro.Product.ID = pid
		pro.Product.Type = "configure"
		pro.Product.Batch_price = &min
		x := min * (float64(1) + *secondpriceratio)
		pro.Product.Second_price = &x
		SQL := util.InsertSQLBuild(pro.Product, "products", []string{"Is_sale", "Created_at", "Updated_at", "Deleted", "Recommended"})
		if _, err := r.db.NamedExec(SQL, pro.Product); err != nil {
			return nil, err
		}
		// 插入图片
		ImageSQL := `INSERT INTO product_images (id, product_id, big_image) values ($1,$2,$3)`
		if pro.Images != nil {
			for _, image := range *pro.Images {
				_, err = r.db.Exec(ImageSQL, xid.New().String(), pro.Product.ID, image)
				if err != nil {
					return nil, err
				}
			}
		}
		//循环插入子商品信息
		productname := pro.Product.Name
		for _, v := range *pro.Spec {
			// 插入子商品
			pro.Product.ID = xid.New().String()
			pro.Product.Parent_id = &pid
			if v.Spec_2 == nil {
				pro.Product.Name = productname + " " + v.Spec_1
			} else {
				pro.Product.Name = productname + " " + v.Spec_1 + " " + *v.Spec_2
			}
			pro.Product.Spec_1 = &v.Spec_1
			pro.Product.Spec_2 = v.Spec_2
			pro.Product.Batch_price = &v.Batch_price
			pro.Product.Type = "simple"
			x := v.Batch_price * (float64(1) + *secondpriceratio)
			pro.Product.Second_price = &x
			SQL := util.InsertSQLBuild(pro.Product, "products", []string{"Is_sale", "Created_at", "Updated_at", "Deleted", "Recommended"})
			if _, err := r.db.NamedExec(SQL, pro.Product); err != nil {
				return nil, err
			}
			// 插入商家子商品关联表
			RelSQL := `INSERT INTO rel_merchants_products (product_id,merchant_id,stock,retail_price) values ($1,$2,$3,$4)`
			RetailPrice := *pro.Product.Batch_price * (float64(1) + *retailpriceratio)
			_, err = r.db.Exec(RelSQL, pro.Product.ID, pro.Product.Merchant_id, v.Stock, RetailPrice)
			if err != nil {
				return nil, err
			}
			// 插入子商品图片
			ImageSQL := `INSERT INTO product_images (id, product_id, big_image) values ($1,$2,$3)`
			if v.Images != nil {
				for _, image := range *v.Images {
					_, err = r.db.Exec(ImageSQL, xid.New().String(), pro.Product.ID, image)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	product, err := r.FindByID(pid)
	if err != nil {
		return nil, err
	}
	return product, nil
}

//更新商品
func (r *ProductRepository) Update(ctx context.Context, pro *model.ProductUpdateARR) (*model.Product, error) {
	// 查找品牌,获取溢价比率
	brand, err := L("brand").(*BrandRepository).FindByID(ctx, pro.Product.Brand_id)
	if err != nil {
		return nil, err
	}
	secondpriceratio := &brand.Second_price_ratio
	retailpriceratio := &brand.Retail_price_ratio
	// 获取当前用户信息
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "provider" {
		return nil, errors.New("只有供货商才能更新商品")
	}
	a := pro.Product.ID
	pro.Product.Merchant_id = *gc.CurrentUser(ctx)
	productType, parent_id, err := r.FindProductType(ctx, pro.Product.ID)
	if err != nil {
		return nil, err
	}

	if *productType == "simple" {
		// 如果规格为空，直接更新父商品
		if pro.Spec == nil {
			if pro.BatchPrice == nil || pro.Stock == nil {
				return nil, errors.New("价格或者库存不能为空")
			}
			pro.Product.Batch_price = pro.BatchPrice
			x := *pro.BatchPrice * (float64(1) + *secondpriceratio)
			pro.Product.Second_price = &x
			// 更新商品表
			SQL := util.UpdateSQLBuild(pro.Product, "products", []string{"Is_sale", "Created_at", "Updated_at", "Deleted", "Recommended", "Spec_1_name", "Spec_2_name", "Spec_1", "Spec_2", "Parent_id", "Type"})
			if _, err := r.db.NamedExec(SQL, pro.Product); err != nil {
				return nil, err
			}
			//更新商家商品关联表
			RelSQL := `update rel_merchants_products set stock=$1,retail_price=$2,origin_price=(CASE WHEN origin_price IS NULL THEN NULL ELSE cast( $2 as numeric) END) where product_id=$3 and merchant_id=$4`
			RetailPrice := *pro.BatchPrice * (float64(1) + *retailpriceratio)
			_, err = r.db.Exec(RelSQL, pro.Stock, RetailPrice, pro.Product.ID, pro.Product.Merchant_id)
			if err != nil {
				return nil, err
			}
			//先删除商品图片再插入
			DelSQL := "DELETE from product_images where product_id=$1"
			if _, err := r.db.Exec(DelSQL, pro.Product.ID); err != nil {
				return nil, err
			}
			ImageSQL := `INSERT INTO product_images (id, product_id, big_image) values ($1,$2,$3)`
			for _, image := range *pro.Images {
				if _, err = r.db.Exec(ImageSQL, xid.New().String(), pro.Product.ID, image); err != nil {
					return nil, err
				}
			}
		} else {
			// 这这种情况是第一次上传没有开启规格，随后又开启
			if parent_id != nil {
				return nil, errors.New("子商品不能再添加规格")
			}
			if pro.BatchPrice == nil || pro.Stock == nil {
				return nil, errors.New("价格或者库存不能为空")
			}
			//找最小价格
			min := (*pro.Spec)[0].Batch_price
			for _, v := range *pro.Spec {
				if min > v.Batch_price {
					min = v.Batch_price
				}
			}
			pro.Product.Batch_price = &min
			x := min * (float64(1) + *secondpriceratio)
			pro.Product.Second_price = &x
			pro.Product.Type = "configure"
			// 更新商品表
			SQL := util.UpdateSQLBuild(pro.Product, "products", []string{"Is_sale", "Created_at", "Updated_at", "Deleted", "Recommended", "Parent_id"})
			if _, err := r.db.NamedExec(SQL, pro.Product); err != nil {
				return nil, err
			}
			//更新商家商品关联表
			RelSQL := `update rel_merchants_products set stock=$1,retail_price=$2,origin_price=(CASE WHEN origin_price IS NULL THEN NULL ELSE cast( $2 as numeric) END) where product_id=$3 and merchant_id=$4`
			RetailPrice := *pro.BatchPrice * (float64(1) + *retailpriceratio)
			_, err = r.db.Exec(RelSQL, pro.Stock, RetailPrice, pro.Product.ID, pro.Product.Merchant_id)
			if err != nil {
				return nil, err
			}
			//先删除商品图片再插入
			DelSQL := "DELETE from product_images where product_id=$1"
			if _, err := r.db.Exec(DelSQL, pro.Product.ID); err != nil {
				return nil, err
			}
			ImageSQL := `INSERT INTO product_images (id, product_id, big_image) values ($1,$2,$3)`
			for _, image := range *pro.Images {
				if _, err = r.db.Exec(ImageSQL, xid.New().String(), pro.Product.ID, image); err != nil {
					return nil, err
				}
			}

			// 添加规格商品
			pid := pro.Product.ID
			productname := pro.Product.Name
			for _, v := range *pro.Spec {
				// 插入子商品
				pro.Product.ID = xid.New().String()
				pro.Product.Parent_id = &pid
				if v.Spec_2 == nil {
					pro.Product.Name = productname + " " + v.Spec_1
				} else {
					pro.Product.Name = productname + " " + v.Spec_1 + " " + *v.Spec_2
				}
				pro.Product.Spec_1 = &v.Spec_1
				pro.Product.Spec_2 = v.Spec_2
				pro.Product.Batch_price = &v.Batch_price
				pro.Product.Type = "simple"
				x := v.Batch_price * (float64(1) + *secondpriceratio)
				pro.Product.Second_price = &x
				SQL := util.InsertSQLBuild(pro.Product, "products", []string{"Is_sale", "Created_at", "Updated_at", "Deleted", "Recommended"})
				if _, err := r.db.NamedExec(SQL, pro.Product); err != nil {
					return nil, err
				}
				// 插入商家子商品关联表
				RelSQL := `INSERT INTO rel_merchants_products (product_id,merchant_id,stock,retail_price) values ($1,$2,$3,$4)`
				RetailPrice := *pro.Product.Batch_price * (float64(1) + *retailpriceratio)
				_, err = r.db.Exec(RelSQL, pro.Product.ID, pro.Product.Merchant_id, v.Stock, RetailPrice)
				if err != nil {
					return nil, err
				}
				// 插入子商品图片
				ImageSQL := `INSERT INTO product_images (id, product_id, big_image) values ($1,$2,$3)`
				if v.Images != nil {
					for _, image := range *v.Images {
						_, err = r.db.Exec(ImageSQL, xid.New().String(), pro.Product.ID, image)
						if err != nil {
							return nil, err
						}
					}
				}
			}

		}
	} else {
		// 更新父商品
		SQL := util.UpdateSQLBuild(pro.Product, "products", []string{"Is_sale", "Created_at", "Updated_at", "Deleted", "Recommended", "Type", "Batch_price", "Second_price"})
		if _, err := r.db.NamedExec(SQL, pro.Product); err != nil {
			return nil, err
		}

		//先删除父商品图片再插入
		DelSQL := "DELETE from product_images where product_id=$1"
		if _, err := r.db.Exec(DelSQL, pro.Product.ID); err != nil {
			return nil, err
		}
		ImageSQL := `INSERT INTO product_images (id, product_id, big_image) values ($1,$2,$3)`
		for _, image := range *pro.Images {
			if _, err = r.db.Exec(ImageSQL, xid.New().String(), pro.Product.ID, image); err != nil {
				return nil, err
			}
		}
		// 更新子商品
		productname := pro.Product.Name
		pid := pro.Product.ID
		for _, v := range *pro.Spec {
			if v.Product_id == nil {
				// 添加了规格就直接插入子商品
				proid := xid.New().String()
				pro.Product.ID = proid
				pro.Product.Parent_id = &pid
				if v.Spec_2 == nil {
					pro.Product.Name = productname + " " + v.Spec_1
				} else {
					pro.Product.Name = productname + " " + v.Spec_1 + " " + *v.Spec_2
				}
				pro.Product.Spec_1 = &v.Spec_1
				pro.Product.Spec_2 = v.Spec_2
				pro.Product.Batch_price = &v.Batch_price
				pro.Product.Type = "simple"
				x := v.Batch_price * (float64(1) + *secondpriceratio)
				pro.Product.Second_price = &x
				SQL := util.InsertSQLBuild(pro.Product, "products", []string{"Is_sale", "Created_at", "Updated_at", "Deleted", "Recommended"})
				if _, err := r.db.NamedExec(SQL, pro.Product); err != nil {
					return nil, err
				}

				// 插入商家子商品关联表
				SQL2 := `INSERT INTO rel_merchants_products (product_id,merchant_id,stock,retail_price) values ($1,$2,$3,$4)`
				RetailPrice := *pro.Product.Batch_price * (float64(1) + *retailpriceratio)
				_, err = r.db.Exec(SQL2, proid, pro.Product.Merchant_id, v.Stock, RetailPrice)
				if err != nil {
					return nil, err
				}

				// 插入子商品图片
				ImageSQL := `INSERT INTO product_images (id, product_id, big_image) values ($1,$2,$3)`
				if v.Images != nil {
					for _, image := range *v.Images {
						_, err = r.db.Exec(ImageSQL, xid.New().String(), proid, image)
						if err != nil {
							return nil, err
						}
					}
				}
			} else {
				//更新了规格就更新子商品
				pro.Product.ID = *v.Product_id
				if v.Spec_2 == nil {
					pro.Product.Name = productname + " " + v.Spec_1
					// 因为下面UpdateSQLBuild函数会忽略nil字段，之前可能设置了规格2，这里没有了，所以重新设置为null
					nullSQL := `UPDATE products set spec_2 = null,spec_2_name = null where id = $1 OR parent_id = $2 `
					if _, err := r.db.Exec(nullSQL, a, a); err != nil {
						return nil, err
					}
				} else {
					pro.Product.Name = productname + " " + v.Spec_1 + " " + *v.Spec_2
				}
				pro.Product.Spec_1 = &v.Spec_1
				pro.Product.Spec_2 = v.Spec_2
				pro.Product.Batch_price = &v.Batch_price
				pro.Product.Type = "simple"
				x := v.Batch_price * (float64(1) + *secondpriceratio)
				pro.Product.Second_price = &x
				SQL := util.UpdateSQLBuild(pro.Product, "products", []string{"Is_sale", "Created_at", "Updated_at", "Deleted", "Recommended", "Parent_id"})
				if _, err := r.db.NamedExec(SQL, pro.Product); err != nil {
					return nil, err
				}

				//更新商家商品关联表
				RelSQL := `update rel_merchants_products set stock=$1,retail_price=$2,origin_price=(CASE WHEN origin_price IS NULL THEN NULL ELSE cast( $2 as numeric) END) where product_id=$3 and merchant_id=$4`
				RetailPrice := v.Batch_price * (float64(1) + *retailpriceratio)
				_, err = r.db.Exec(RelSQL, v.Stock, RetailPrice, v.Product_id, pro.Product.Merchant_id)
				if err != nil {
					return nil, err
				}
				//先删除商品图片再插入
				DelSQL := "DELETE from product_images where product_id=$1"
				if _, err := r.db.Exec(DelSQL, v.Product_id); err != nil {
					return nil, err
				}
				ImageSQL := `INSERT INTO product_images (id, product_id, big_image) values ($1,$2,$3)`
				if v.Images != nil {
					for _, image := range *v.Images {
						if _, err = r.db.Exec(ImageSQL, xid.New().String(), v.Product_id, image); err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}
	product, err := r.FindByID(a)
	if err != nil {
		return nil, err
	}
	return product, nil
}

//通过ID查找商品
func (p *ProductRepository) FindByID(ID string) (*model.Product, error) {
	product := &model.Product{}
	SQL := `SELECT * FROM products WHERE id = $1`
	row := p.db.QueryRowx(SQL, ID)
	err := row.StructScan(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// 通过ID集合查找商品
func (p *ProductRepository) FindByIDs(ids []string) ([]*model.Product, error) {
	products := make([]*model.Product, 0)
	SQL := util.NewSQLBuilder(products).BuildQuery() + ` WHERE id in (?)`
	query, args, err := sqlx.In(SQL, ids)
	if err != nil {
		return nil, err
	}
	query = p.db.Rebind(query)
	if err = p.db.Select(&products, query, args...); err != nil {
		return nil, err
	}
	return products, nil
}

//查找子商品
func (p *ProductRepository) FindChildren(ID string) ([]*model.Product, error) {
	products := []*model.Product{}
	SQL := "SELECT * FROM products WHERE parent_id = $1 AND deleted = false"
	if err := p.db.Select(&products, SQL, ID); err != nil {
		return nil, err
	}
	return products, nil
}

//根据条件查找商品
func (p *ProductRepository) Search(ctx context.Context, first *int32, offset *int32, search *model.ProductSearchInput, sort *model.ProductSortInput) ([]*model.Product, error) {
	var SQL string
	products := make([]*model.Product, 0)
	fetchSize := util.GetInt32(first, defaultListFetchSize)
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if offset != nil {
		if *usertype == "admin" {
			builder := util.NewSQLBuilder(products)
			p.searchWhere(search, builder)
			SQL = builder.BuildQuery() + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		} else if *usertype == "provider" {
			builder := util.NewSQLBuilder(products)
			p.searchWhere(search, builder)
			SQL = builder.BuildQuery() + gc.WhereMerchant(ctx) + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		} else if *usertype == "ally" {
			builder := util.NewSQLBuilder(nil)
			p.searchWhere(search, builder)
			where := builder.BuildWhere()
			SQL = `SELECT pro.* FROM rel_merchants_products rmp LEFT JOIN products pro ON rmp.product_id=pro.id ` + where + ` AND rmp.merchant_id = '` + *gc.CurrentUser(ctx) + `'` + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		} else {
			builder := util.NewSQLBuilder(nil)
			p.searchWhere(search, builder)
			where := builder.BuildWhere()
			SQL = `SELECT pro.* FROM rel_agents_products rmp LEFT JOIN products pro ON rmp.product_id=pro.id ` + where + ` AND rmp.agent_id = '` + *gc.CurrentUser(ctx) + `'` + p.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		}
		if err := p.db.Select(&products, SQL, fetchSize, *offset); err != nil {
			return nil, err
		}
		return products, nil
	} else {
		if *usertype == "admin" {
			builder := util.NewSQLBuilder(products)
			p.searchWhere(search, builder)
			SQL = builder.BuildQuery() + p.buildSortSQL(sort) + ` LIMIT $1;`
		} else if *usertype == "provider" {
			builder := util.NewSQLBuilder(products)
			p.searchWhere(search, builder)
			SQL = builder.BuildQuery() + gc.WhereMerchant(ctx) + p.buildSortSQL(sort) + ` LIMIT $1;`
		} else if *usertype == "ally" {
			builder := util.NewSQLBuilder(nil)
			p.searchWhere(search, builder)
			where := builder.BuildWhere()
			SQL = `SELECT pro.* FROM rel_merchants_products rmp LEFT JOIN products pro ON rmp.product_id=pro.id ` + where + ` AND rmp.merchant_id = '` + *gc.CurrentUser(ctx) + `'` + p.buildSortSQL(sort) + ` LIMIT $1;`
		} else {
			builder := util.NewSQLBuilder(nil)
			p.searchWhere(search, builder)
			where := builder.BuildWhere()
			SQL = `SELECT pro.* FROM rel_agents_products rmp LEFT JOIN products pro ON rmp.product_id=pro.id ` + where + ` AND rmp.agent_id = '` + *gc.CurrentUser(ctx) + `'` + p.buildSortSQL(sort) + ` LIMIT $1;`
		}
		if err := p.db.Select(&products, SQL, fetchSize); err != nil {
			return nil, err
		}
		return products, nil
	}

}

func (p *ProductRepository) buildSortSQL(sort *model.ProductSortInput) (s string) {
	if sort == nil {
		return
	}
	s = " ORDER BY " + sort.OrderBy + " " + util.GetString(sort.Sort, "ASC") + " "
	return
}

func (p *ProductRepository) searchWhere(search *model.ProductSearchInput, builder *util.SQLBuilder) {
	builder.WhereStruct(search, true).WhereRow("deleted=false AND parent_id is NULL")
	if search != nil && search.Key != nil {
		builder.WhereRow(fmt.Sprintf("name ILIKE '%%%s%%'", *search.Key))
	}
}

// 根据不同角色和搜索条件获取商品数量
func (p *ProductRepository) Count(ctx context.Context, c *model.ProductSearchInput) (int, error) {
	var count int
	var SQL string
	builder := util.NewSQLBuilder(nil)
	p.searchWhere(c, builder)
	where := builder.BuildWhere()
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return 0, err
	}
	if *status != "normal" {
		return 0, errors.New("用户状态不正常")
	}
	if *usertype == "ally" {
		SQL = `SELECT count(*) FROM rel_merchants_products rmp LEFT JOIN products pro ON rmp.product_id=pro.id ` + where + ` AND rmp.merchant_id = '` + *gc.CurrentUser(ctx) + `'`
	}
	if *usertype == "admin" {
		SQL = `SELECT count(*) FROM products ` + where
	}
	if *usertype == "provider" {
		SQL = `SELECT count(*) FROM products ` + where + gc.WhereMerchant(ctx)
	}
	if *usertype == "agent" {
		SQL = `SELECT count(*) FROM rel_agents_products rmp LEFT JOIN products pro ON rmp.product_id=pro.id ` + where + ` AND rmp.agent_id = '` + *gc.CurrentUser(ctx) + `'`
	}
	if err := p.db.Get(&count, SQL); err != nil {
		return 0, err
	}
	return count, nil
}

// 根据商品ID获取图片信息
func (p *ProductRepository) FindImagesById(ID string) ([]*model.ProductImage, error) {
	images := []*model.ProductImage{}
	SQL := "SELECT * FROM product_images WHERE product_id = $1 ORDER BY created_at"
	if err := p.db.Select(&images, SQL, ID); err != nil {
		return nil, err
	}
	return images, nil
}

//根据商品ID获取供货商（上传商品用户）资料
func (p *ProductRepository) FindMerchants(productID string) (*model.MerchantProfile, error) {

	profiles := &model.MerchantProfile{}

	SQL := `SELECT * FROM merchant_profiles mp WHERE mp.user_id=(SELECT merchant_id FROM products WHERE id=$1)`
	row := p.db.QueryRowx(SQL, productID)
	err := row.StructScan(profiles)
	if err != nil {
		return nil, err
	}
	profiles.CompanyAddress, err = model.NewAddress(profiles.CompanyAddressRow)
	if err != nil {
		return nil, err
	}
	profiles.DeliveryAddress, err = model.NewAddress(profiles.DeliveryAddressRow)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

// 根据商品ID查找当前商家售价
func (r *ProductRepository) FindRetailPrice(ctx context.Context, ID string) (*float64, error) {
	var price float64
	SQL := `SELECT retail_price FROM rel_merchants_products where product_id=$1 and merchant_id=$2`
	err := r.db.Get(&price, SQL, ID, gc.CurrentUser(ctx))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &price, nil
}

// 根据商品ID查找售价供admin查看（即供货商）
func (r *ProductRepository) FindRetailPriceForAdmin(ctx context.Context, ID string) (*float64, error) {
	var price float64
	SQL := `SELECT rmp.retail_price FROM rel_merchants_products rmp left join products pro on pro.merchant_id=rmp.merchant_id where rmp.product_id=$1`
	err := r.db.Get(&price, SQL, ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &price, nil
}

// 查找库存信息
func (r *ProductRepository) GetStock(ctx context.Context, ID string) (*int32, error) {
	var stock int32
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "ally" || *usertype == "provider" {
		SQL := `select stock from rel_merchants_products where product_id=$1 and merchant_id=$2`
		err := r.db.Get(&stock, SQL, ID, gc.CurrentUser(ctx))
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &stock, nil
	} else {
		var SQL string
		if *usertype == "agent" {
			SQL = `SELECT rmp.stock FROM rel_agents_products rap LEFT JOIN rel_merchants_products rmp on rap.product_id=rmp.product_id LEFT JOIN users u ON u."id"=rmp.merchant_id WHERE u."type"='provider' AND rap.product_id=$1`
		} else {
			SQL = `select stock from rel_merchants_products rmp LEFT JOIN products pro on rmp.merchant_id=pro.merchant_id where product_id=$1`
		}
		err := r.db.Get(&stock, SQL, ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &stock, nil
	}
}

// 查找销量信息
func (r *ProductRepository) GetSalesVolume(ctx context.Context, ID string) (*int32, error) {
	var SalesVolume int32
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "ally" || *usertype == "provider" {
		SQL := `select sales_volume from rel_merchants_products where product_id=$1 and merchant_id=$2`
		err := r.db.Get(&SalesVolume, SQL, ID, gc.CurrentUser(ctx))
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &SalesVolume, nil
	} else if *usertype == "agent" {
		SQL := `select sales_volume from rel_agents_products where product_id=$1 and agent_id=$2`
		err := r.db.Get(&SalesVolume, SQL, ID, gc.CurrentUser(ctx))
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &SalesVolume, nil
	} else {
		SQL := `SELECT tt.sale FROM (SELECT SUM(sales_volume) sale,product_id FROM rel_merchants_products GROUP BY product_id) tt WHERE tt.product_id=$1`
		err := r.db.Get(&SalesVolume, SQL, ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &SalesVolume, nil
	}
}

// 查找收藏信息
func (r *ProductRepository) GetFavorites(ctx context.Context, ID string) (*int32, error) {
	var favorites int32
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "ally" || *usertype == "provider" || *usertype == "agent" {
		SQL := `select count(*) from favorites where product_id=$1 and merchant_id=$2`
		err := r.db.Get(&favorites, SQL, ID, gc.CurrentUser(ctx))
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &favorites, nil
	} else {
		SQL := `select count(*) from favorites where product_id=$1`
		err := r.db.Get(&favorites, SQL, ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &favorites, nil
	}
}

//获取商品类型和父ID
func (r *ProductRepository) FindProductType(ctx context.Context, ID string) (*string, *string, error) {
	var res struct {
		Type      string
		Parent_id *string
	}
	SQL := `select type, parent_id from products where id=$1`
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(&res); err != nil {
		return nil, nil, err
	}
	return &res.Type, res.Parent_id, nil
}

// 删除商品
func (r *ProductRepository) Deleted(ctx context.Context, ID string) (bool, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}

	// 加盟商在售不能删除
	var count int
	countSQL := `SELECT COUNT(*) FROM rel_merchants_products rmp LEFT JOIN users u ON rmp.merchant_id=u."id" WHERE rmp.product_id = $1 AND rmp.stock>0 AND u.type='ally'`
	if err := r.db.Get(&count, countSQL, ID); err != nil {
		return false, err
	}
	if count > 0 {
		return false, errors.New("此商品有加盟商进货在售，不能删除")
	}

	// 有子商品不能删除
	var childrencount int
	childrencountSQL := `SELECT COUNT(*) FROM products WHERE parent_id = $1 AND deleted=false`
	if err := r.db.Get(&childrencount, childrencountSQL, ID); err != nil {
		return false, err
	}
	if childrencount > 0 {
		return false, errors.New("此商品有子商品未被删除，请删除子商品后重试")
	}

	if *usertype == "ally" || *usertype == "agent" {
		return false, errors.New("您不能删除商品")
	} else if *usertype == "provider" {
		SQL = `update products set deleted=true where id=$1` + gc.WhereMerchant(ctx)
	} else {
		SQL = `update products set deleted=true where id=$1`
	}

	if _, err := r.db.Exec(SQL, ID); err != nil {
		return false, err
	}
	return true, nil
}

// 加盟商设置自己商品的售价
func (r *ProductRepository) SetReatilPriceByAlly(ctx context.Context, price float64, ID string) (*model.Product, error) {
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "ally" {
		return nil, errors.New("您不能执行此操作")
	}
	SQL = `update rel_merchants_products set retail_price = $1,origin_price=(CASE WHEN origin_price IS NULL THEN NULL ELSE cast( $1 as numeric) END) where product_id = $2` + gc.WhereMerchant(ctx)

	result, err := r.db.Exec(SQL, price, ID)
	if err != nil {
		return nil, err
	}

	if af, err := result.RowsAffected(); err != nil {
		return nil, err
	} else if af != 1 {
		return nil, fmt.Errorf("设置失败, 检查ID为 %s 的商品是否存在", ID)
	}

	product, err := r.FindByID(ID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// 商品上下架
func (r *ProductRepository) SetIsSale(ctx context.Context, ID string) (bool, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}

	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}

	if *usertype == "admin" {
		return false, errors.New("管理员不能进行此设置")
	}

	// 只有Type为simple的商品才能上下架
	product, err := r.FindByID(ID)
	if err != nil {
		return false, err
	}
	if product.Type == "configure" {
		return false, errors.New("此商品不能上下架")
	}

	var SQL string
	if *usertype == "agent" {
		SQL = `update rel_agents_products set is_sale = not is_sale where product_id = $1` + gc.WhereMerchant(ctx)
	} else {
		SQL = `update rel_merchants_products set is_sale = not is_sale where product_id = $1` + gc.WhereMerchant(ctx)
	}

	result, err := r.db.Exec(SQL, ID)
	if err != nil {
		return false, err
	}

	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("设置失败, 检查ID为 %s 的商品是否存在", ID)
	}

	return true, nil
}

// 当前会话用户根据商品ID查看商品是否上架
func (r *ProductRepository) CatIsSale(ctx context.Context, ID string) (*bool, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}

	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}

	// if *usertype == "admin" {
	// 	return nil, errors.New("只有供货商或者加盟商能查看自己商品上下架信息")
	// }

	var isSale bool
	var SQL string
	if *usertype == "agent" {
		SQL = `SELECT is_sale FROM rel_agents_products WHERE product_id = $1` + gc.WhereMerchant(ctx)
	} else {
		SQL = `SELECT is_sale FROM rel_merchants_products WHERE product_id = $1` + gc.WhereMerchant(ctx)
	}

	if err := r.db.Get(&isSale, SQL, ID); err != nil {
		return nil, err
	}

	return &isSale, nil
}
