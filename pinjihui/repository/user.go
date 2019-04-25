package repository

import (
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"

    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "github.com/rs/xid"
    "golang.org/x/net/context"
    valid "gopkg.in/asaskevich/govalidator.v9"
    gc "pinjihui.com/pinjihui/context"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/util"
)

const (
    DefaultListFetchSize = 10
    PLATFORM             = "platform"
)

type UserRepository struct {
    BaseRepository
}

func NewUserRepository(db *sqlx.DB, log *logging.Logger) *UserRepository {
    return &UserRepository{BaseRepository{db: db, log: log}}
}

func (u *UserRepository) FindBy(column, value string) (*model.User, error) {
    user := &model.User{}

    userSQL := util.NewSQLBuilder(user).WhereRow(column + "=$1").BuildQuery()
    row := u.db.QueryRowx(userSQL, value)
    err := row.StructScan(user)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    if err != nil {
        u.log.Errorf("Error in retrieving user : %v", err)
        return nil, err
    }
    return user, nil
}

func (u *UserRepository) FindByMobile(mobile string) (*model.User, error) {
    return u.FindBy("mobile", mobile)
}

func (u *UserRepository) FindByID(ID string) (*model.User, error) {
    return u.FindBy("id", ID)
}

func (u *UserRepository) FindByInviteCode(InviteCode string) (*model.User, error) {
    return u.FindBy("invite_code", InviteCode)
}

func (u *UserRepository) CreateUser(user *model.User) (*model.User, error) {
    valid.TagMap["hasAlphaAndNumeric"] = valid.Validator(func(str string) bool {
        return regexp.MustCompile("[0-9]").MatchString(str) &&
            regexp.MustCompile("[a-zA-Z]").MatchString(str)

    })
    _, err := valid.ValidateStruct(user)
    if err != nil {
        return nil, err
    }
    userId := xid.New()
    user.ID = userId.String()
    err = u.Save(user)
    if err != nil {
        return nil, err
    }
    return u.FindByID(userId.String())
}

func (u *UserRepository) Save(user *model.User) error {
    userSQL := `INSERT INTO users (id, mobile, name, password, last_ip, type) VALUES (:id, :mobile, :name, :password, :last_ip, :type)`
    user.HashedPassword()
    _, err := u.db.NamedExec(userSQL, user)
    return err
}

func (u *UserRepository) ComparePassword(userCredentials *model.UserCredentials) (*model.User, error) {
    user, err := u.FindByMobile(userCredentials.Mobile)
    if err != nil {
        return nil, errors.New(gc.UserNameOrPasswordError)
    }
    if result := user.ComparePassword(userCredentials.Password); !result {
        return nil, errors.New(gc.UserNameOrPasswordError)
    }
    return user, nil
}

func (u *UserRepository) FindWechart(openid string) (*model.WechartUser, error) {
    wxUser := model.WechartUser{}
    SQLS := `SELECT openid, session_key, user_id FROM wecharts WHERE openid=$1`
    row := u.db.QueryRowx(SQLS, openid)
    err := row.StructScan(&wxUser)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    if err != nil {
        u.log.Errorf("Error in retrieving WechartUser : %v", err)
        return nil, err
    }
    return &wxUser, nil
}

func (u *UserRepository) FindWechartByUserID(id string) (*model.WechartUser, error) {
    wxUser := model.WechartUser{}
    SQLS := `SELECT openid, session_key, user_id FROM wecharts WHERE user_id=$1`
    row := u.db.QueryRowx(SQLS, id)
    err := row.StructScan(&wxUser)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    if err != nil {
        u.log.Errorf("Error in retrieving WechartUser : %v", err)
        return nil, err
    }
    return &wxUser, nil
}

func (u *UserRepository) CreateWxUser(ctx context.Context, openid string, sessionKey string, tx *sqlx.Tx) (*model.WechartUser, error) {
    wxUser := &model.WechartUser{
        OpenId:     openid,
        SessionKey: sessionKey,
        NickName:   model.DefaultWXNickName,
    }
    userID, err := u.getOrCreateUserIfNotAuthorized(ctx, wxUser)
    if err != nil {
        return nil, err
    }
    wxUser.UserID = userID
    SQL := `INSERT INTO wecharts(openid, session_key, nick_name, gender, user_id) VALUES(:openid, :session_key, :nick_name, :gender, :user_id)`

    _, err = tx.NamedExec(SQL, wxUser)
    if err != nil {
        return nil, err
    }
    return wxUser, nil
}

func (u *UserRepository) getOrCreateUserIfNotAuthorized(ctx context.Context, wxUser *model.WechartUser) (string, error) {
    if ctx.Value("is_authorized").(bool) {
        return *gc.CurrentUser(ctx), nil
    }
    user := model.User{
        ID:       xid.New().String(),
        Password: wxUser.SessionKey,
        Name:     &wxUser.NickName,
        LastIp:   *ctx.Value("requester_ip").(*string),
        Type:     model.UTConsumer,
    }
    return user.ID, u.Save(&user)
}

func (u *UserRepository) WxLogin(ctx context.Context, code string) (*model.WechartUser, error) {
    resp, err := u.requestOpenid(ctx, code)
    if err != nil {
        return nil, err
    }
    wxUser, err := u.FindWechart(resp.Openid)
    if err == gc.ErrNoRecord {
        //创建记录
        tx := u.db.MustBegin()
        wxUser, err = u.CreateWxUser(ctx, resp.Openid, resp.Skey, tx)
        //送新人优惠券,无须填邀请码
        if err != nil {
            return nil, err
        }
        err = u.addCouponForUser(tx, model.ForFirstLogin, wxUser.UserID)
        if err != nil {
            tx.Rollback()
            return nil, err
        }
        tx.Commit()
    }
    return wxUser, err
}

type wxLoginResp struct {
    Openid string
    Skey   string
}

func (u *UserRepository) requestOpenid(ctx context.Context, code string) (*wxLoginResp, error) {
    appID := ctx.Value("config").(*gc.Config).WechartAppID
    secret := ctx.Value("config").(*gc.Config).WechartSecret
    url := fmt.Sprintf(`https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code`, appID, secret, code)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    var respObj interface{}
    err = json.Unmarshal(body, &respObj)
    if err != nil {
        return nil, err
    }
    u.log.Info("wechart login request return: %+v", respObj)
    respMap := respObj.(map[string]interface{})
    if _, ok := respMap["errcode"]; ok {
        return nil, errors.New(respMap["errmsg"].(string))
    }
    openID := respMap["openid"].(string)
    skey := respMap["session_key"].(string)
    return &wxLoginResp{Openid: openID, Skey: skey}, nil
}

func (u *UserRepository) UpdateLoginStatus(ctx context.Context, id string) error {
    query := `UPDATE users SET last_ip=$1, last_login_time=CURRENT_TIMESTAMP WHERE id=$2`
    _, err := u.db.Exec(query, ctx.Value("requester_ip").(*string), id)
    if err != nil {
        u.log.Errorf("Update login status failed: %v", err)
    }
    return err
}

func (u *UserRepository) ReceiveCoupon(ctx context.Context, inviteCode string) (bool, error) {
    user, err := u.FindByID(*gc.CurrentUser(ctx))
    if err != nil {
        return false, err
    }
    if user.Invited == true {
        return false, errors.New("您已领取过了，不能重复领取")
    } else {
        var count, newercount, invitercount int
        // 输入邀请码是否存在
        SQL := `SELECT count(*) FROM users where invite_code = $1`
        if err := u.db.Get(&count, SQL, inviteCode); err != nil {
            return false, err
        }
        // 新人优惠券数量
        newerSQL := `SELECT count(*) FROM coupons WHERE type = 'for_newer'`
        if err := u.db.Get(&newercount, newerSQL); err != nil {
            return false, err
        }
        // 邀请人优惠券数量
        inviterSQL := `SELECT count(*) FROM coupons WHERE type = 'for_inviter'`
        if err := u.db.Get(&invitercount, inviterSQL); err != nil {
            return false, err
        }

        if count <= 0 {
            return false, errors.New("邀请码不存在")
        } else {
            if inviteCode == user.InviteCode {
                return false, errors.New("不能用自己的邀请码领取优惠券")
            }
            // 开启事务
            tx, err := u.db.Begin()
            if err != nil {
                return false, err
            }
            // 新人领取优惠券
            //SQL := `INSERT into user_coupons SELECT $1, $2, description,value,FALSE,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,expired_at,limit_amount,'for_newer', start_at,merchant_id  FROM coupons WHERE "type" = 'for_newer'`
            for a := 0; a < newercount; a++ {
                SQL := `INSERT into user_coupons SELECT
                           $1,
                           $2,
                           description,
                           value,
                           FALSE,
                           CURRENT_TIMESTAMP,
                           CURRENT_TIMESTAMP,
                           case when validity_days is null
                             then expired_at
                           else current_date + validity_days end as expired_at,
                           limit_amount,
                           'for_newer',
                           case when validity_days is null
                             then start_at
                           else current_date end                 AS start_at,
                           merchant_id
                         FROM coupons
                         WHERE "type" = 'for_newer'
                         ORDER BY created_at ASC
                         LIMIT 1
                         OFFSET $3`
                result, err := tx.Exec(SQL, xid.New().String(), gc.CurrentUser(ctx), a)
                if err != nil {
                    tx.Rollback()
                    return false, err
                }
                if af, err := result.RowsAffected(); err != nil {
                    tx.Rollback()
                    return false, err
                } else if af != 1 {
                    tx.Rollback()
                    return false, errors.New("领取失败")
                }
            }

            // 更新状态
            userSQL := `UPDATE users SET invited = true WHERE id = $1`
            result2, err2 := tx.Exec(userSQL, gc.CurrentUser(ctx))
            if err2 != nil {
                tx.Rollback()
                return false, err2
            }
            if af, err := result2.RowsAffected(); err != nil {
                tx.Rollback()
                return false, err
            } else if af != 1 {
                tx.Rollback()
                return false, errors.New("更新失败")
            }

            // 邀请人领取优惠券
            for b := 0; b < invitercount; b++ {
                SQLSTR := `INSERT into user_coupons SELECT $1, $2, description,value,FALSE,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,case when validity_days is null
                             then expired_at
                           else current_date + validity_days end as expired_at,limit_amount,'for_inviter', case when validity_days is null
                             then start_at
                           else current_date end                 AS start_at,merchant_id  FROM coupons WHERE "type" = 'for_inviter' ORDER BY created_at ASC LIMIT 1 OFFSET $3`
                userbycode, errbycode := u.FindByInviteCode(inviteCode)
                if errbycode != nil {
                    tx.Rollback()
                    return false, errbycode
                }
                result3, err3 := tx.Exec(SQLSTR, xid.New().String(), userbycode.ID, b)
                if err3 != nil {
                    tx.Rollback()
                    return false, err3
                }
                if af, err := result3.RowsAffected(); err != nil {
                    tx.Rollback()
                    return false, err
                } else if af != 1 {
                    tx.Rollback()
                    return false, errors.New("领取失败")
                }
            }

            tx.Commit()
            return true, nil
        }
    }
}

type WxUser struct {
    NickName  string  `db:"nick_name"`
    Gender    int32
    Language  *string
    City      *string
    Province  *string
    Country   *string
    AvatarUrl *string `db:"avatar_url"`
    UserID    *string `db:"user_id"`
}

func (u *UserRepository) SaveWxUserInfo(ctx context.Context, user *WxUser) (bool, error) {
    gc.CheckAuth(ctx)
    user.UserID = gc.CurrentUser(ctx)
    query := `UPDATE wecharts SET nick_name=:nick_name,gender=:gender,language=:language,city=:city,province=:province, country=:country,avatar_url=:avatar_url WHERE user_id=:user_id`
    _, err := u.db.NamedExec(query, user)
    return true, err
}

func (u *UserRepository) AddShareCoupon(ctx context.Context, couponType string) error {
    gc.CheckAuth(ctx)
    userId := *gc.CurrentUser(ctx)
    has, err := u.HasShareCoupon(ctx, couponType)
    if err != nil {
        return err
    }
    if !has {
        tx := u.db.MustBegin()
        //为分享者增加优惠券
        err = u.addCouponForUser(tx, couponType, userId)
        if err != nil {
            return err
        }
        tx.Commit()
    }
    return nil
}

func (u *UserRepository) addCouponForUser(tx *sqlx.Tx, couponType, userId string) error {
    query := `INSERT into user_coupons SELECT $1, $2, description,value,FALSE,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,case when validity_days is null
        then expired_at
        else current_date + validity_days end as expired_at,limit_amount,type, case when validity_days is null
        then start_at
        else current_date end AS start_at,merchant_id  FROM coupons WHERE "type" = $3 ORDER BY created_at ASC LIMIT 1`
    _, err := tx.Exec(query, xid.New().String(), userId, couponType)
    return err
}

func (u *UserRepository) HasShareCoupon(ctx context.Context, couponType string) (bool, error) {
    gc.CheckAuth(ctx)
    //分享者是否存在可用的分享优惠券
    query := fmt.Sprintf(`SELECT COUNT(*) FROM user_coupons WHERE user_id=$1 AND expired_at >= CURRENT_DATE AND start_at<=CURRENT_DATE AND type='%s' AND used=FALSE`, couponType)
    sharerCouponCount, err := u.GetIntFormDB(query, gc.CurrentUser(ctx))
    return sharerCouponCount > 0, err
}

func (u *UserRepository) BindPhoneNumber(ctx context.Context, number string) error {
   gc.CheckAuth(ctx)
   if !(valid.IsNumeric(number) && valid.StringLength(number, "11", "11")) {
       return errors.New("mobile invalid")
   }
   query := `UPDATE users SET mobile=$1 WHERE id=$2`
   _, err := u.db.Exec(query, number, gc.CurrentUser(ctx))
   return err
}
