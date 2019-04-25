package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/util"
)

type RegionRepository struct {
    BaseRepository
}

func NewRegionRepository(db *sqlx.DB, log *logging.Logger) *RegionRepository {
    return &RegionRepository{BaseRepository{db: db, log: log}}
}

func (r *RegionRepository) FindAllParents(id int32) ([]*model.Region, error) {
    regions := make([]*model.Region, 0)

    sql := `WITH RECURSIVE region_tree(id, parent_id, name, sort_order) AS (
                SELECT * FROM regions WHERE id = $1
                UNION ALL SELECT r.* FROM regions r, region_tree t 
                WHERE r.id = t.parent_id
            )
            select * from region_tree;`
    if err := r.db.Select(&regions, sql, id); err != nil{
        return nil, err
    }
    return regions, nil
}

func (r *RegionRepository) FindByParentID(parent *int32) ([]*model.Region, error) {
    regions := make([]*model.Region, 0)
    sql := `SELECT * FROM regions WHERE parent_id=$1`
    if err := r.db.Select(&regions, sql, util.GetInt32(parent, 0)); err != nil {
        return nil, err
    }
    return regions, nil
}
