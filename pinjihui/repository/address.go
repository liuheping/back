package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
    "errors"
    "fmt"
    "pinjihui.com/pinjihui/util"
    "github.com/rs/xid"
    "golang.org/x/net/context"
    gc "pinjihui.com/pinjihui/context"
    valid "gopkg.in/asaskevich/govalidator.v9"
    "database/sql"
)

type AddressRepository struct {
    BaseRepository
}

const ADDR_TABLE = "addresses"

func NewAddressRepository(db *sqlx.DB, log *logging.Logger) *AddressRepository {
    return &AddressRepository{BaseRepository{db: db, log: log}}
}

func CheckAC(ctx context.Context, address *model.ShippingAddress) {
    if *gc.CurrentUser(ctx) != address.UserId {
        panic(gc.ErrUnAC)
    }
}

func (r *AddressRepository) FindByID(ctx context.Context, ID string) (*model.ShippingAddress, error) {
    address, err := r.FindOriginByID(ctx, ID)
    if err != nil {
        return nil, err
    }
    if address.ShippingAddress.Address, err = model.NewAddress(&address.Address); err != nil {
        return nil, err
    }
    return &address.ShippingAddress, nil
}

func (r *AddressRepository) FindOriginByID(ctx context.Context, id string) (*model.ShippingAddressDB, error) {
    gc.CheckAuth(ctx)
    address := &model.ShippingAddressDB{}

    addressSQL := `SELECT id,user_id,consignee,address,zipcode,mobile,is_default FROM addresses WHERE id = $1`
    row := r.db.QueryRowx(addressSQL, id)
    err := row.StructScan(address)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    if err != nil {
        return nil, err
    }
    CheckAC(ctx, &address.ShippingAddress)
    return address, nil
}

//如果没有记录,则既没有error, 也没有数据,即返回(nil,nil)
func (r *AddressRepository) FindAll(ctx context.Context) ([]*model.ShippingAddress, error) {
    gc.CheckAuth(ctx)

    addressSQL := `SELECT id,user_id,consignee,address,zipcode,mobile,is_default FROM addresses WHERE user_id = $1`
    rows, err := r.db.Queryx(addressSQL, *gc.CurrentUser(ctx))
    if err != nil {
        return nil, err
    }
    addrs := []*model.ShippingAddress{}
    for rows.Next() {
        address := model.ShippingAddressDB{}
        err = rows.StructScan(&address)
        if err != nil {
            return nil, err
        }
        if address.ShippingAddress.Address, err = model.NewAddress(&address.Address); err != nil {
            return nil, err
        }
        addrs = append(addrs, &address.ShippingAddress)
    }
    return addrs, nil
}

func (r *AddressRepository) Save(ctx context.Context, newAddr *model.ShippingAddress) (*model.ShippingAddress, error) {
    gc.CheckAuth(ctx)
    //填充省市区id和省市区全名
    _, err := valid.ValidateStruct(newAddr)
    if err != nil {
        return nil, err
    }
    if err := r.FillAddr(newAddr.Address); err != nil {
        return nil, err
    }
    var addrSQL string
    if newAddr.ID == "" {
        //create
        newAddr.UserId = *gc.CurrentUser(ctx)
        newAddr.ID = xid.New().String()
        addrSQL = util.NewSQLBuilder(newAddr).Table("addresses").InsertSQLBuild(nil)
    } else {
        if _, err := r.FindByID(ctx, newAddr.ID); err != nil {
            return nil, err
        }
        addrSQL = util.NewSQLBuilder(newAddr).Table("addresses").UpdateSQLBuild([]string{"UserId", "IsDefault"})
    }
    if _, err := r.db.NamedExec(addrSQL, newAddr); err != nil {
        return nil, err
    }
    return newAddr, nil
}

func (r *AddressRepository) SetDefault(ctx context.Context, id string) (bool, error) {
    gc.CheckAuth(ctx)
    sqls := `UPDATE addresses SET is_default=(id=$1) WHERE user_id=$2`
    if result, err := r.db.Exec(sqls, id, *gc.CurrentUser(ctx)); err != nil {
        return false, err
    } else if af, _ := result.RowsAffected(); af < 1 {
        return false, fmt.Errorf("Update failed, check if address exist, id: %s", id)
    }
    return true, nil
}

func (r *AddressRepository) FillAddr(address *model.Address) error {
    if address == nil || address.AreaId == 0 {
        return errors.New("invalid area")
    }
    regions, err := L("region").(*RegionRepository).FindAllParents(address.AreaId)
    if err != nil {
        return err
    }
    if len(regions) != 3 {
        return errors.New("invalid area")
    }
    address.CityId = regions[1].ID
    address.ProvinceId = regions[2].ID
    rname := fmt.Sprintf("%s %s %s", regions[2].Name, regions[1].Name, regions[0].Name)
    address.RegionName = &rname
    return nil
}

func (r *AddressRepository) Delete(ctx context.Context, id string) (bool, error) {
    gc.CheckAuth(ctx)
    sqls := `DELETE FROM addresses WHERE id=$1 AND user_id=$2`
    result, err := r.db.Exec(sqls, id, *gc.CurrentUser(ctx))
    if err != nil {
        return false, err
    }
    if af, err := result.RowsAffected(); err != nil {
        return false, err
    } else if af != 1 {
        return false, fmt.Errorf("Update failed, check if address exist, id: %s", id)
    }
    return true, nil
}
