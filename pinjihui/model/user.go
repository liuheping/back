package model

import (
    "log"
    "golang.org/x/crypto/bcrypt"
)

const (
    UTConsumer                = "consumer"
    UTAgent                   = "agent"
    DefaultWXNickName         = "微信用户"
    DefaultMobileUserNickName = "微信用户"
)

type User struct {
    ID            string
    Mobile        *string `valid:"numeric,required,length(11|11)~mobile invalid"`
    Email         *string
    Name          *string
    Password      string  `valid:"required,length(8|100)~password_length_error,printableascii~password_printable_ascii_error,hasAlphaAndNumeric~password_must_strong"`
    LastIp        string  `db:"last_ip"`
    CreatedAt     string  `db:"created_at"`
    UpdatedAt     string  `db:"updated_at"`
    Type          string
    Status        string
    LastLoginTime *string `db:"last_login_time"`
    InviteCode    string  `db:"invite_code"`
    Invited       bool
    Openid        string  `fi:"-"`
}

func (user *User) HashedPassword() error {
    hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Println(err)
        return err
    }
    user.Password = string(hash)
    return nil
}

func (user *User) ComparePassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        log.Println(err)
        return false
    }
    return true
}

type WechartUser struct {
    OpenId     string
    SessionKey string `db:"session_key"`
    UserID     string `db:"user_id"`
    NickName   string `db:"nick_name"`
    Gender     int16
}
