package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/util"
)

const (
	defaultListFetchSize = 10
)

type UserRepository struct {
	db             *sqlx.DB
	roleRepository *RoleRepository
	log            *logging.Logger
}

func NewUserRepository(db *sqlx.DB, roleRepository *RoleRepository, log *logging.Logger) *UserRepository {
	return &UserRepository{db: db, roleRepository: roleRepository, log: log}
}

func (u *UserRepository) FindByMobile(mobile string) (*model.User, error) {
	user := &model.User{}
	SQL := `SELECT * FROM users WHERE mobile = $1 AND type != 'consumer'`
	row := u.db.QueryRowx(SQL, mobile)
	if err := row.StructScan(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	user.ID = xid.New().String()
	var count int
	countSQL := `select count(1) from users where mobile=$1`
	if err := u.db.Get(&count, countSQL, user.Mobile); err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("The phone number has already existed")
	}
	SQL := `INSERT INTO users (id, type, password, last_ip, mobile, status) VALUES (:id, :type, :password, :last_ip, :mobile, :status)`
	user.HashedPassword()
	user.Status = "unchecked"
	if _, err := u.db.NamedExec(SQL, user); err != nil {
		return nil, err
	}
	return user, nil
}

//根据条件查找用户
func (u *UserRepository) Search(ctx context.Context, first *int32, offset *int32, search *model.UserSearchInput, sort *model.UserSortInput) ([]*model.User, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return nil, errors.New("您不能查看用户列表")
	}
	users := make([]*model.User, 0)
	fetchSize := util.GetInt32(first, defaultListFetchSize)
	builder := util.NewSQLBuilder(users)
	u.searchWhere(search, builder)
	if offset != nil {
		SQL := builder.BuildQuery() + u.buildSortSQL(sort) + ` LIMIT $1 OFFSET $2;`
		if err := u.db.Select(&users, SQL, fetchSize, *offset); err != nil {
			return nil, err
		}
		return users, nil
	}
	SQL := builder.BuildQuery() + u.buildSortSQL(sort) + ` LIMIT $1;`
	if err := u.db.Select(&users, SQL, fetchSize); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserRepository) buildSortSQL(sort *model.UserSortInput) (s string) {
	if sort == nil {
		return
	}
	s = " ORDER BY " + sort.OrderBy + " " + util.GetString(sort.Sort, "ASC") + " "
	return
}

func (u *UserRepository) searchWhere(search *model.UserSearchInput, builder *util.SQLBuilder) {
	builder.WhereStruct(search, true)
	if search != nil && search.Name != nil {
		builder.WhereRow(fmt.Sprintf("name ILIKE '%%%s%%'", *search.Name))
	}
}

func (u *UserRepository) Count(ctx context.Context, c *model.UserSearchInput) (int, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return 0, err
	}
	if *status != "normal" {
		return 0, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return 0, errors.New("您不能查看用户列表")
	}
	var count int
	builder := util.NewSQLBuilder(nil)
	u.searchWhere(c, builder)
	where := builder.BuildWhere()
	SQL := `SELECT count(*) FROM users ` + where
	if err := u.db.Get(&count, SQL); err != nil {
		return 0, err
	}
	return count, nil
}

func (u *UserRepository) ComparePassword(userCredentials *model.UserCredentials) (*model.User, error) {
	user, err := u.FindByMobile(userCredentials.Mobile)
	if err != nil {
		return nil, err
	}
	if result := user.ComparePassword(userCredentials.Password); !result {
		return nil, errors.New(gc.UnauthorizedAccess)
	}
	return user, nil
}

func (u *UserRepository) FindByID(ID string) (*model.User, error) {
	user := &model.User{}
	SQL := `SELECT * FROM users WHERE id = $1`
	row := u.db.QueryRowx(SQL, ID)
	if err := row.StructScan(user); err != nil {
		return nil, err
	}
	return user, nil
}

//UpdatePasswordByID 通过ID修改密码
func (u *UserRepository) UpdatePasswordByID(NewPassword string, ID string) (bool, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	hashpwd := string(hash)
	SQL := `update users set password=$1 where id=$2`
	if _, err := u.db.Exec(SQL, hashpwd, ID); err != nil {
		return false, err
	}
	return true, nil
}

//UpdatePasswordByMobileAndType 找回密码
func (u *UserRepository) UpdatePasswordByMobileAndType(Mobile string, NewPassword string) (bool, error) {
	var count int
	SQL := `select count(1) from users where mobile=$1`
	if err := u.db.Get(&count, SQL, Mobile); err != nil {
		return false, err
	}
	if count != 1 {
		return false, errors.New("账户异常")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	hashpwd := string(hash)
	updateSQL := `update users set password=$1 where mobile=$2 AND type != 'consumer'`
	if _, err := u.db.Exec(updateSQL, hashpwd, Mobile); err != nil {
		return false, err
	}
	return true, nil
}

//修改商户资料，没有直接添加
func (u *UserRepository) UpdateProfileByID(ctx context.Context, Profile *model.MerchantProfile) (*model.MerchantProfile, error) {
	var count int
	var SQL string
	usertype, _, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	countSQL := `select count(*) from merchant_profiles where user_id=$1`
	if err := u.db.Get(&count, countSQL, Profile.UserId); err != nil {
		return nil, err
	}
	if err := L("address").(*AddressRepository).fillAddr(Profile.CompanyAddress); err != nil {
		return nil, err
	}
	if err := L("address").(*AddressRepository).fillAddr(Profile.DeliveryAddress); err != nil {
		return nil, err
	}
	if count <= 0 {
		// 插入时，填写的区域ID不能和已有的加盟商重复
		if *usertype == "ally" {
			if _, err := u.CheckAreaIdForAlly(Profile.CompanyAddress.AreaId); err != nil {
				return nil, err
			}
		}
		SQL = util.InsertSQLBuild(Profile, "merchant_profiles", []string{"Balance", "CompanyAddressRow", "DeliveryAddressRow", "Created_at", "Updated_at"})
	} else {
		// 更新时，加盟商不能修改公司地址的区域ID
		if *usertype == "ally" {
			mp, err := u.MerchantProfile(Profile.UserId)
			if err != nil {
				return nil, err
			}
			if mp.CompanyAddress.AreaId != Profile.CompanyAddress.AreaId {
				return nil, errors.New("加盟商公司地址所在的区域不能修改")
			}
		}
		// 更新后设置用户状态为未审核(管理员除外)
		if *usertype != "admin" {
			if _, err := u.SetUserStatus(Profile.UserId); err != nil {
				return nil, err
			}
		}
		SQL = util.UpdateSQLBuild2(Profile, "merchant_profiles", []string{"Balance", "CompanyAddressRow", "DeliveryAddressRow", "Created_at", "Updated_at"})
	}
	if _, err := u.db.NamedExec(SQL, Profile); err != nil {
		return nil, err
	}
	return Profile, nil
}

//获取提款资料列表
func (u *UserRepository) TakeCashList(UserID string) (*[]*model.TakeCash, error) {
	takecash := new([]*model.TakeCash)
	SQL := `SELECT * FROM debit_card_info where user_id=$1`
	err := u.db.Select(takecash, SQL, UserID)
	if err != nil {
		return nil, err
	}
	for _, v := range *takecash {
		v.DebitCardInfo, err = model.NewDebitCardInfo(v.DebitCardInfoRow)
		if err != nil {
			return nil, err
		}
	}
	return takecash, nil
}

//修改提款资料,没有就添加一条，有就直接修改
func (u *UserRepository) UpdateTakeCash(DebitCard *model.TakeCash) (*model.TakeCash, error) {
	var count int
	countSQL := `select count(1) from debit_card_info where user_id=$1`
	err := u.db.Get(&count, countSQL, DebitCard.UserID)
	if err != nil {
		return nil, err
	}
	var SQL string
	if count <= 0 {
		SQL = util.InsertSQLBuild(DebitCard, "debit_card_info", []string{"IsChecked", "Updated_at", "Created_at", "DebitCardInfoRow"})
	} else {
		SQL = util.UpdateSQLBuild2(DebitCard, "debit_card_info", []string{"UserID", "IsChecked", "Updated_at", "Created_at", "DebitCardInfoRow"})
	}
	if _, err := u.db.NamedExec(SQL, DebitCard); err != nil {
		return nil, err
	}
	return DebitCard, nil
}

//根据UserID获取商户资料
func (u *UserRepository) MerchantProfile(UserID string) (*model.MerchantProfile, error) {
	profiles := &model.MerchantProfile{}
	SQL := `SELECT * FROM merchant_profiles where user_id=$1`
	row := u.db.QueryRowx(SQL, UserID)
	err := row.StructScan(profiles)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
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

//获取地区资料
func (u *UserRepository) Region(ParentID *int32) ([]*model.Region, error) {
	r := make([]*model.Region, 0)
	SQL := `SELECT * FROM regions where parent_id=$1`

	err := u.db.Select(&r, SQL, ParentID)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// 审核商家（资料）
func (u *UserRepository) CheckMerchant(ctx context.Context, ID string) (bool, error) {
	usertype, status, err := L("public").(*PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return false, err
	}
	if *status != "normal" {
		return false, errors.New("用户状态不正常")
	}
	if *usertype != "admin" {
		return false, errors.New("只有管理员能审核商户资料")
	}
	SQL := `UPDATE users SET status = $1 WHERE id=$2`
	if _, err := u.db.Exec(SQL, "normal", ID); err != nil {
		return false, err
	}
	return true, nil
}

// 检查加盟商区域ID是否重复(添加时)
func (u *UserRepository) CheckAreaIdForAlly(AreaID int32) (bool, error) {
	SQL := `SELECT COUNT(*) FROM merchant_profiles mp LEFT JOIN users u ON mp.user_id = u."id" WHERE u."type"= 'provider' AND  (mp.company_address).area_id = $1`
	var count int
	err := u.db.Get(&count, SQL, AreaID)
	if err != nil {
		return false, err
	}
	if count >= 1 {
		return false, errors.New("同一区域只能有一个加盟商")
	} else {
		return true, nil
	}
}

// 设置用户状态
func (u *UserRepository) SetUserStatus(ID string) (bool, error) {
	SQL := `UPDATE users SET status = $1 WHERE id = $2`
	result, err := u.db.Exec(SQL, "unchecked", ID)
	if err != nil {
		return false, err
	}
	if af, err := result.RowsAffected(); err != nil {
		return false, err
	} else if af != 1 {
		return false, fmt.Errorf("更新失败, 检查商家ID是否存在: %s", ID)
	}
	return true, nil
}
