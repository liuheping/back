package service

import (
	"encoding/base64"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/op/go-logging"
	xctx "golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
)

type AuthService struct {
	appName *string
	// signedSecret        *string
	expiredTimeInSecond *time.Duration
	log                 *logging.Logger
}

func NewAuthService(config *context.Config, log *logging.Logger) *AuthService {
	//return &AuthService{&config.AppName, &config.JWTSecret, &config.JWTExpireIn, log}
	return &AuthService{&config.AppName, &config.JWTExpireIn, log}
}

//SignJWT 签名JWT
func (a *AuthService) SignJWT(user *model.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         base64.StdEncoding.EncodeToString([]byte(user.ID)),
		"created_at": user.CreatedAt,
		"exp":        time.Now().Add(time.Second * *a.expiredTimeInSecond).Unix(),
		"iss":        *a.appName,
	})

	tokenString, err := token.SignedString([]byte(user.Password))
	return &tokenString, err
}

//ValidateJWT 验证JWT
func (a *AuthService) ValidateJWT(ctx xctx.Context, tokenString *string) (*jwt.Token, error) {
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("	unexpected signing method: %v", token.Header["alg"])
		}
		//return []byte(*a.signedSecret), nil

		userID, err := base64.StdEncoding.DecodeString(token.Claims.(jwt.MapClaims)["id"].(string))
		if err != nil {
			return nil, err
		}

		user, err := ctx.Value("userRepository").(*repository.UserRepository).FindByID(string(userID))
		if err != nil || user == nil {
			return nil, err
		}
		return []byte(user.Password), nil

	})
	return token, err
}
