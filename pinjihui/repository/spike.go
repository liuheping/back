package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/util"
    "database/sql"
    "qiniupkg.com/x/log.v7"
    "errors"
    gc "pinjihui.com/pinjihui/context"
    "fmt"
)

type SpikeRepository struct {
    BaseRepository
}

func NewSpikeRepository(db *sqlx.DB, log *logging.Logger) *SpikeRepository {
    return &SpikeRepository{BaseRepository{db: db, log: log}}
}

func (s *SpikeRepository) List(first int, offset int) ([]*model.Spike, error) {
    spikes := make([]*model.Spike, 0)
    builder := util.NewSQLBuilder(spikes).WhereRow("spikes.expired_at > CURRENT_TIMESTAMP").
        Join("products", "spikes.product_id=products.id").
        WhereRow("products.is_sale=true AND products.deleted=false")
    sqlStr := builder.OrderBy("start_at", "ASC").Limit(first, &offset).BuildQuery()
    if err := s.db.Select(&spikes, sqlStr, builder.Args...); err != nil {
        return nil, err
    }
    return spikes, nil
}

func (s *SpikeRepository) Count() (int, error) {
    var count int
    slqs := `SELECT count(*) FROM spikes Join products ON spikes.product_id=products.id WHERE spikes.expired_at > CURRENT_TIMESTAMP AND products.is_sale=true AND products.deleted=false`
    err := s.db.Get(&count, slqs)
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (s *SpikeRepository) UpdateProductPrice() {
    sqls := `UPDATE rel_merchants_products r
SET  origin_price=retail_price,retail_price = price from spikes s
                      WHERE
                        s.product_id = r.product_id AND r.merchant_id = s.merchant_id AND start_at < current_timestamp
                        and expired_at > current_timestamp AND total_count!=0 AND retail_price!=price;
                        UPDATE rel_merchants_products r
SET retail_price = origin_price FROM spikes s
WHERE
  s.product_id = r.product_id AND r.merchant_id = s.merchant_id AND retail_price!=origin_price AND origin_price IS NOT NULL AND expired_at = (select max(expired_at)
                                                                                  from spikes s2
                                                                                  where s2.product_id = r.product_id AND
                                                                                        r.merchant_id = s2.merchant_id)
  AND (expired_at < current_timestamp OR total_count = 0);`
    if _, err := s.db.Exec(sqls); err != nil {
        s.log.Errorf("update product spike price failed, err:%v", err)
    }
}

func (s *SpikeRepository) IsSpiking(productId string) (bool, error) {
    query := `SELECT COUNT(*) FROM spikes WHERE product_id=$1 AND start_at<current_timestamp AND expired_at>current_timestamp AND total_count > 0`
    c, err := s.GetIntFormDB(query, productId)
    return c != 0, err
}

func (s *SpikeRepository) CanSpike(tx *sqlx.Tx, productId, merchantId string, productQuantity int) (bool, error) {
    var countLimit struct {
        Buy_limit   int
        Total_count int
    }
    query := `SELECT total_count,buy_limit FROM spikes WHERE product_id=$1 AND merchant_id=$2 AND start_at<current_timestamp AND expired_at>current_timestamp`
    err := tx.Get(&countLimit, query, productId, merchantId)
    //没有参加秒杀
    if err == sql.ErrNoRows {
        return false, nil
    }

    if err != nil {
        log.Errorf("get spikes count failed: %v", err)
        return false, err
    }
    //已秒光,普通商品
    if countLimit.Total_count == 0 {
        return false, nil
    }
    if countLimit.Total_count < productQuantity {
        return false, errors.New("没有足够的秒杀商品")
    }
    if countLimit.Buy_limit < productQuantity {
        return false, fmt.Errorf("购物车中秒杀商品每人限购%d件", countLimit.Buy_limit)
    }
    return true, nil
}

func (s *SpikeRepository) FindSpikeByPM(productID, merchantID string) (*model.Spike, error) {
    spike := model.Spike{}
    query := util.NewSQLBuilder(&spike).WhereRow("product_id=$1 AND merchant_id=$2").BuildQuery()
    err := s.db.Get(&spike, query, productID, merchantID)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    return &spike, err
}
