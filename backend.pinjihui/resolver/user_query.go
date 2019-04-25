package resolver

import (
	"errors"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/op/go-logging"
	"golang.org/x/net/context"
	gcontext "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/loader"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

func (r *Resolver) User(ctx context.Context, args struct {
	ID string
}) (*userResolver, error) {
	userId := ctx.Value("user_id").(*string)
	user, err := loader.LoadUser(ctx, args.ID)
	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
		return nil, err
	}
	ctx.Value("log").(*logging.Logger).Debugf("Retrieved user by user_id[%s] : %v", *userId, *user)
	return &userResolver{user}, nil
}

func (r *Resolver) Login(ctx context.Context, args struct {
	Mobile   string
	Password string
}) (*userResolver, error) {
	userCredentials := &model.UserCredentials{
		Mobile:   args.Mobile,
		Password: args.Password,
	}
	user, err := ctx.Value("userRepository").(*repository.UserRepository).ComparePassword(userCredentials)
	if err != nil {
		return nil, err
	}
	if user.Type == "consumer" {
		return nil, errors.New("消费者不能登陆管理后台")
	}
	return &userResolver{user}, nil
}

func (r *Resolver) TakeCashs(ctx context.Context) (*[]*takeCashResolver, error) {

	if isAuthorized := ctx.Value("is_authorized").(bool); !isAuthorized {
		return nil, errors.New(gcontext.CredentialsError)
	}
	userId := ctx.Value("user_id").(*string)

	takecashprofile, err := ctx.Value("userRepository").(*repository.UserRepository).TakeCashList(*userId)

	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
		return nil, err
	}
	ctx.Value("log").(*logging.Logger).Debugf("Retrieved takecashprofile by user_id[%s] : %v", userId, takecashprofile)

	l := make([]*takeCashResolver, len(*takecashprofile))
	for i := range l {
		l[i] = &takeCashResolver{(*takecashprofile)[i]}
	}
	return &l, nil
}

//获取当前会话商户资料
func (r *Resolver) MerchantProfiles(ctx context.Context) (*merchantProfileResolver, error) {
	if isAuthorized := ctx.Value("is_authorized").(bool); !isAuthorized {
		return nil, errors.New(gcontext.CredentialsError)
	}
	userId := ctx.Value("user_id").(*string)
	mp, err := ctx.Value("userRepository").(*repository.UserRepository).MerchantProfile(*userId)
	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
		return nil, err
	}
	if mp == nil {
		return nil, nil
	}
	ctx.Value("log").(*logging.Logger).Debugf("Retrieved merchantprofile by user_id[%s] : %v", userId, mp)
	return &merchantProfileResolver{mp}, nil
}

//获取地区列表
func (r *Resolver) Regions(ctx context.Context, args struct{ Pid *int32 }) (*[]*regoinResolver, error) {

	rg, err := ctx.Value("userRepository").(*repository.UserRepository).Region(args.Pid)

	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
		return nil, err
	}
	ctx.Value("log").(*logging.Logger).Debugf("Retrieved rgion : %v", rg)

	l := make([]*regoinResolver, len(rg))
	for i := range l {
		l[i] = &regoinResolver{(rg)[i]}
	}
	return &l, nil
}

//根据地址ID获取用户收货地址
func (r *Resolver) MyAddress(ctx context.Context, args struct {
	ID graphql.ID
}) (*shippingAddressResolver, error) {
	ID := string(args.ID)
	addr, err := rp.L("address").(*rp.AddressRepository).FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &shippingAddressResolver{addr}, nil
}

//获取会话用户所有收货地址
func (r *Resolver) MyAllAddress(ctx context.Context) (*[]*shippingAddressResolver, error) {
	if isAuthorized := ctx.Value("is_authorized").(bool); !isAuthorized {
		return nil, errors.New(gcontext.CredentialsError)
	}
	userId := ctx.Value("user_id").(*string)
	addr, err := rp.L("address").(*rp.AddressRepository).FindAll(ctx, *userId)
	if err != nil {
		return nil, err
	}
	l := make([]*shippingAddressResolver, len(addr))
	for i := range l {
		l[i] = &shippingAddressResolver{(addr)[i]}
	}
	return &l, nil
}

//根据条件获取所有用户
func (r *Resolver) Users(ctx context.Context, args struct {
	First  *int32
	Offset *int32
	Search *model.UserSearchInput
	Sort   *model.UserSortInput
}) (*usersConnectionResolver, error) {
	fetchSize := util.GetInt32(args.First, DefaultPageSize)
	list, err := ctx.Value("userRepository").(*repository.UserRepository).Search(ctx, &fetchSize, args.Offset, args.Search, args.Sort)
	if err != nil {
		return nil, err
	}
	count, err := ctx.Value("userRepository").(*repository.UserRepository).Count(ctx, args.Search)
	if err != nil {
		return nil, err
	}
	var from *string
	var to *string
	if len(list) > 0 {
		from = &(list[0].ID)
		to = &(list[len(list)-1].ID)
	}
	return &usersConnectionResolver{m: list, totalCount: count, from: from, to: to, hasNext: util.If(len(list) == int(fetchSize), true, false).(bool)}, nil
}

func (r *Resolver) Me(ctx context.Context) (*userResolver, error) {
	if isAuthorized := ctx.Value("is_authorized").(bool); !isAuthorized {
		return nil, errors.New(gcontext.CredentialsError)
	}
	userId := ctx.Value("user_id").(*string)
	user, err := ctx.Value("userRepository").(*repository.UserRepository).FindByID(*userId)
	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
		return nil, err
	}
	ctx.Value("log").(*logging.Logger).Debugf("Retrieved user by user_id[%s] : %v", *userId, *user)
	return &userResolver{user}, nil
}
