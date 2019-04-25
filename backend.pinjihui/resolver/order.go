package resolver

import (
	"database/sql"
	"errors"
	"time"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
	rp "pinjihui.com/backend.pinjihui/repository"
)

type orderResolver struct {
	m *model.Order
}

func (r *orderResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *orderResolver) User(ctx context.Context) (*userResolver, error) {
	user, err := ctx.Value("userRepository").(*repository.UserRepository).FindByID(r.m.User_id)
	if err != nil {
		return nil, err
	}
	return &userResolver{user}, nil
}

func (r *orderResolver) ChildrenOrders(ctx context.Context) (*[]*orderResolver, error) {
	orders, err := rp.L("order").(*rp.OrderRepository).FindChildren(ctx, r.m.ID)
	if err != nil {
		return nil, err
	}
	l := make([]*orderResolver, len(orders))
	for i := range l {
		l[i] = &orderResolver{(orders)[i]}
	}
	return &l, nil
}

func (r *orderResolver) Address() (*orderAddressResolver, error) {
	return &orderAddressResolver{r.m.Address}, nil
}

func (r *orderResolver) Merchant() (*merchantProfileResolver, error) {
	profile, err := rp.L("order").(*rp.OrderRepository).FindMerchant(r.m.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &merchantProfileResolver{profile}, nil
}

func (r *orderResolver) Products() (*[]*productInOrderResolver, error) {
	pros, err := rp.L("productinorder").(*rp.ProductInOrderRepository).FindByOrderID(r.m.ID)
	if err != nil {
		return nil, err
	}
	l := make([]*productInOrderResolver, len(pros))
	for i := range l {
		l[i] = &productInOrderResolver{(pros)[i]}
	}
	return &l, nil
}

func (r *orderResolver) ShippingId() *graphql.ID {
	if r.m.Shipping_id == nil {
		return nil
	}
	res := graphql.ID(*r.m.Shipping_id)
	return &res
}

func (r *orderResolver) ShippingName() *string {
	return r.m.Shipping_name
}

func (r *orderResolver) PayCode() graphql.ID {
	return graphql.ID(r.m.Pay_code)
}

func (r *orderResolver) PayName() string {
	return r.m.Pay_name
}

func (r *orderResolver) Amount() float64 {
	return r.m.Amount
}

func (r *orderResolver) ShippingFee() *float64 {
	return r.m.Shipping_fee
}

func (r *orderResolver) MoneyPaid() float64 {
	return *r.m.Money_paid
}

func (r *orderResolver) OrderAmount() float64 {
	return *r.m.Order_amount
}

func (r *orderResolver) CreatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Created_at)
	return graphql.Time{Time: res}, err
}

func (r *orderResolver) ConfirmTime() (*graphql.Time, error) {
	if r.m.Confirm_time == nil {
		return nil, nil
	}
	res, err := time.Parse(time.RFC3339, *r.m.Confirm_time)
	return &graphql.Time{Time: res}, err
}

func (r *orderResolver) PayTime() (*graphql.Time, error) {
	if r.m.Pay_time == nil {
		return nil, nil
	}
	res, err := time.Parse(time.RFC3339, *r.m.Pay_time)
	return &graphql.Time{Time: res}, err
}

func (r *orderResolver) ShippingTime() (*graphql.Time, error) {
	if r.m.Shipping_time == nil {
		return nil, nil
	}
	res, err := time.Parse(time.RFC3339, *r.m.Shipping_time)
	return &graphql.Time{Time: res}, err
}

func (r *orderResolver) Note() *string {
	return r.m.Note
}

func (r *orderResolver) Postscript() *string {
	return r.m.Postscript
}

func (r *orderResolver) UsedCoupon() *string {
	return r.m.Used_coupon
}

func (r *orderResolver) Status() string {
	return r.m.Status
}

func (r *orderResolver) OfferAmount() float64 {
	return r.m.Offer_amount
}

func (r *orderResolver) Income(ctx context.Context) (*incomeResolver, error) {
	var a interface{}
	usertype, status, err := rp.L("public").(*rp.PublicRepository).GetCurrentUserType(ctx)
	if err != nil {
		return nil, err
	}
	if *status != "normal" {
		return nil, errors.New("用户状态不正常")
	}
	if *usertype == "ally" {
		a = &allyIncomeResolver{r.m}
	}
	if *usertype == "provider" {
		a = &providerIncomeResolver{r.m}
	}
	if *usertype == "admin" {
		a = &adminIncomeResolver{r.m}
	}
	return &incomeResolver{a}, nil
}

func (r *orderResolver) Shipping_info(ctx context.Context) (*shippingInfoResolver, error) {
	info, err := rp.L("shippinginfo").(*rp.ShippingInfoRepository).FindByOrderID(ctx, r.m.ID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &shippingInfoResolver{info}, nil
}
