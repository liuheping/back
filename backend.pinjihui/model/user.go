package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string
	Name          *string
	Mobile        *string
	Email         *string
	Password      string
	LastIp        *string `db:"last_ip"`
	CreatedAt     string  `db:"created_at"`
	UpdatedAt     string  `db:"updated_at"`
	LastLoginTime *string `db:"last_login_time"`
	Type          string
	Status        string
	InviteCode    *string `db:"invite_code"`
	Invited       bool
	// TokenString   string
	// Roles []*Role
}

type UserSearchInput struct {
	Name     *string `db:"name"`
	Mobile   *string `db:"mobile"`
	Email    *string `db:"email"`
	Usertype *string `db:"type"`
	Status   *string `db:"status"`
}

type UserSortInput struct {
	OrderBy string
	Sort    *string
}

// 哈希加密密码
func (user *User) HashedPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}
	user.Password = string(hash)
	return nil
}

// 验证密码
func (user *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
