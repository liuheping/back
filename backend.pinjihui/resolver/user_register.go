package resolver

import (
	"errors"

	"github.com/op/go-logging"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) Register(ctx context.Context, args *struct {
	Mobile   string
	Password string
	Type     string
	Code     string
}) (*userResolver, error) {

	user := &model.User{
		Mobile:   &args.Mobile,
		Password: args.Password,
		Type:     args.Type,
		LastIp:   ctx.Value("requester_ip").(*string),
	}

	if args.Code != "1234" {
		ctx.Value("log").(*logging.Logger).Errorf("Message code error")
		return nil, errors.New("message code error,please try again")
	}

	user, err := ctx.Value("userRepository").(*repository.UserRepository).CreateUser(user)

	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
		return nil, err
	}
	ctx.Value("log").(*logging.Logger).Debugf("Created user : %v", *user)

	return &userResolver{user}, nil
}
