package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	valid "gopkg.in/asaskevich/govalidator.v9"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

type AddressRepository struct {
	BaseRepository
}

func NewAddressRepository(db *sqlx.DB, log *logging.Logger) *AddressRepository {
	return &AddressRepository{BaseRepository{db: db, log: log}}
}

func CheckAC(ctx context.Context, address *model.ShippingAddress) {
	if *gc.CurrentUser(ctx) != address.UserId {
		panic(gc.ErrUnAC)
	}
}

func (r *AddressRepository) FindByID(ctx context.Context, ID string) (*model.ShippingAddress, error) {
	gc.CheckAuth(ctx)
	var SQL string
	address := &model.ShippingAddressDB{}
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `SELECT * FROM addresses WHERE id = $1`
	} else {
		SQL = `SELECT * FROM addresses WHERE id = $1` + ` AND user_id = '` + *gc.CurrentUser(ctx) + `'`
	}
	row := r.db.QueryRowx(SQL, ID)
	if err := row.StructScan(address); err != nil {
		return nil, err
	}
	CheckAC(ctx, &address.ShippingAddress)
	if address.ShippingAddress.Address, err = model.NewAddress(address.Address); err != nil {
		return nil, err
	}
	return &address.ShippingAddress, nil
}

//如果没有记录,则既没有error, 也没有数据,即返回(nil,nil)
func (r *AddressRepository) FindAll(ctx context.Context, ID string) ([]*model.ShippingAddress, error) {
	gc.CheckAuth(ctx)
	var SQL string
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "admin" {
		SQL = `SELECT * FROM addresses WHERE user_id = $1`
	} else {
		SQL = `SELECT * FROM addresses WHERE user_id = $1` + ` AND user_id = '` + *gc.CurrentUser(ctx) + `'`
	}
	rows, err := r.db.Queryx(SQL, ID)
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
		if address.ShippingAddress.Address, err = model.NewAddress(address.Address); err != nil {
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
	if err := r.fillAddr(newAddr.Address); err != nil {
		return nil, err
	}
	var SQL string
	if newAddr.ID == "" {
		newAddr.UserId = *gc.CurrentUser(ctx)
		newAddr.ID = xid.New().String()
		SQL = util.InsertSQLBuild(newAddr, "addresses", []string{"Created_at", "Updated_at"})
	} else {
		if _, err := r.FindByID(ctx, newAddr.ID); err != nil {
			return nil, err
		}
		SQL = util.UpdateSQLBuild(newAddr, "addresses", []string{"UserId", "IsDefault", "Created_at", "Updated_at"})
	}
	if _, err := r.db.NamedExec(SQL, newAddr); err != nil {
		return nil, err
	}
	return newAddr, nil
}

func (r *AddressRepository) fillAddr(address *model.Address) error {
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
	// address.CityId = regions[1].ID
	// address.ProvinceId = regions[2].ID
	rname := fmt.Sprintf("%s %s %s", regions[2].Name, regions[1].Name, regions[0].Name)
	address.RegionName = &rname
	return nil
}

// 设置默认收货地址
func (r *AddressRepository) SetDefault(ctx context.Context, id string) (bool, error) {
	gc.CheckAuth(ctx)
	SQL := `UPDATE addresses SET is_default=(id=$1) WHERE user_id=$2`
	if result, err := r.db.Exec(SQL, id, *gc.CurrentUser(ctx)); err != nil {
		return false, err
	} else if af, _ := result.RowsAffected(); af < 1 {
		return false, fmt.Errorf("Update failed, check if address exist, id: %s", id)
	}
	return true, nil
}

// 删除收货地址
func (r *AddressRepository) Delete(ctx context.Context, id string) (bool, error) {
	gc.CheckAuth(ctx)
	SQL := `DELETE FROM addresses WHERE id=$1 AND user_id=$2`
	result, err := r.db.Exec(SQL, id, *gc.CurrentUser(ctx))
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
