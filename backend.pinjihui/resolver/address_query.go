package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
)

func (r *Resolver) Address(ctx context.Context) (*addressResolver, error) {
	//Without using dataloader:
	//user, err := ctx.Value("userRepository").(*repository.userRepository).FindByEmail(args.Email)
	//userId := ctx.Value("user_id").(*string)
	address := &model.Address{}
	//ctx.Value("log").(*logging.Logger).Debugf("Retrieved user by user_id[%s] : %v", *userId, *address)
	return &addressResolver{address}, nil
}
