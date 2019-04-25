package service

import (
	"encoding/base64"
	"fmt"
	"pinjihui.com/pinjihui/context"
	"pinjihui.com/pinjihui/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/op/go-logging"
	"time"
	xctx "golang.org/x/net/context"
	rp "pinjihui.com/pinjihui/repository"
)

type AuthService struct {
	appName             *string
	expiredTimeInSecond *time.Duration
	log                 *logging.Logger
}

func NewAuthService(config *context.Config, log *logging.Logger) *AuthService {
	return &AuthService{&config.AppName, &config.JWTExpireIn, log}
}

func (a *AuthService) SignJWT(user *model.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         base64.StdEncoding.EncodeToString([]byte(user.ID)),
		"created_at": user.CreatedAt,
		"exp":        time.Now().Add(time.Second * *a.expiredTimeInSecond).Unix(),
		"iss":        *a.appName,
		"openid":	  user.Openid,
	})

	tokenString, err := token.SignedString([]byte(user.Password))
	return &tokenString, err
}

func (a *AuthService) ValidateJWT(ctx xctx.Context, tokenString *string) (*jwt.Token, error) {
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("	unexpected signing method: %v", token.Header["alg"])
		}

		userID, err := base64.StdEncoding.DecodeString(token.Claims.(jwt.MapClaims)["id"].(string))
		if err != nil {
		    return nil, err
        }
		user, err := rp.L("user").(*rp.UserRepository).FindByID(string(userID))
		if err != nil || user == nil {
			return nil, err
        }
		return []byte(user.Password), nil
	})
	return token, err
}
