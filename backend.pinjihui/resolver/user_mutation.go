package resolver

import (
	"errors"
	"fmt"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	logging "github.com/op/go-logging"
	"golang.org/x/net/context"
	gcontext "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
)

//RetrievePassword 找回密码
func (r *Resolver) RetrievePassword(ctx context.Context, args *struct {
	Mobile      string
	NewPassword string
	Code        string
}) (bool, error) {
	if args.Code != "1234" {
		ctx.Value("log").(*logging.Logger).Errorf("Message code error")
		return false, errors.New("message code error,please try again")
	}

	issuccessed, err := ctx.Value("userRepository").(*repository.UserRepository).UpdatePasswordByMobileAndType(args.Mobile, args.NewPassword)

	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("update error : %v", err)
		return false, err
	}
	ctx.Value("log").(*logging.Logger).Debugf("is successed: %v", issuccessed)

	return issuccessed, nil
}

//ChangePassword 修改密码
func (r *Resolver) ChangePassword(ctx context.Context, args *struct {
	Oldpassword string
	NewPassword string
}) (bool, error) {

	if isAuthorized := ctx.Value("is_authorized").(bool); !isAuthorized {
		return false, errors.New(gcontext.CredentialsError)
	}
	userId := ctx.Value("user_id").(*string)

	user, err := ctx.Value("userRepository").(*repository.UserRepository).FindByID(*userId)

	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("Find user error : %v", err)
		return false, err
	}

	if result := user.ComparePassword(args.Oldpassword); !result {
		return false, errors.New("原密码错误！")
	}

	issuccessed, err := ctx.Value("userRepository").(*repository.UserRepository).UpdatePasswordByID(args.NewPassword, *userId)

	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("update error : %v", err)
		return false, err
	}

	return issuccessed, nil
}

//更新(添加)商户资料
func (r *Resolver) UpdateProfile(ctx context.Context, args *struct {
	Profile *model.MerchantProfileARR
}) (*merchantProfileResolver, error) {

	if isAuthorized := ctx.Value("is_authorized").(bool); !isAuthorized {
		return nil, errors.New(gcontext.CredentialsError)
	}
	args.Profile.UserId = *ctx.Value("user_id").(*string)

	if args.Profile.CompanyImage != nil {
		str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.Profile.CompanyImage, "\",\""))
		args.Profile.MerchantProfile.CompanyImage = &str
	}

	// if args.Profile.Waiters != nil {
	// 	str := fmt.Sprintf("{\"%s\"}", strings.Join(*args.Profile.Waiters, "\",\""))
	// 	args.Profile.MerchantProfile.Waiters = &str
	// }

	pf, err := ctx.Value("userRepository").(*repository.UserRepository).UpdateProfileByID(ctx, &args.Profile.MerchantProfile)

	if err != nil {
		return nil, err
	}

	return &merchantProfileResolver{pf}, nil
}

//修改提款资料,没有就添加一条，有就直接修改
func (r *Resolver) UpdateTakeCash(ctx context.Context, args *struct {
	Debitcard *model.TakeCash
}) (*takeCashResolver, error) {

	if isAuthorized := ctx.Value("is_authorized").(bool); !isAuthorized {
		return nil, errors.New(gcontext.CredentialsError)
	}
	args.Debitcard.UserID = *ctx.Value("user_id").(*string)

	take_cash := &model.TakeCash{
		UserID:        args.Debitcard.UserID,
		DebitCardInfo: args.Debitcard.DebitCardInfo,
		IsChecked:     false,
	}

	tk, err := ctx.Value("userRepository").(*repository.UserRepository).UpdateTakeCash(take_cash)

	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
		return nil, err
	}
	ctx.Value("log").(*logging.Logger).Debugf("Created take cash : %v", *tk)

	return &takeCashResolver{tk}, nil
}

// 审核商家（资料）
func (r *Resolver) CheckMerchant(ctx context.Context, args struct {
	ID graphql.ID
}) (bool, error) {
	_, err := ctx.Value("userRepository").(*repository.UserRepository).CheckMerchant(ctx, string(args.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}
