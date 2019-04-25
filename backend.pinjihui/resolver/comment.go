package resolver

import (
	"strings"
	"time"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/loader"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

type commentResolver struct {
	m *model.Comment
}

func (r *commentResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *commentResolver) User(ctx context.Context) (*userResolver, error) {
	user, err := loader.LoadUser(ctx, r.m.User_id)
	if err != nil {
		return nil, err
	}
	return &userResolver{user}, nil
}

func (r *commentResolver) Product() (*productResolver, error) {
	product, err := rp.L("product").(*rp.ProductRepository).FindByID(r.m.Product_id)
	if err != nil {
		return nil, err
	}
	return &productResolver{product}, nil
}

func (r *commentResolver) Rank() *int32 {
	return r.m.Rank
}

func (r *commentResolver) Order(ctx context.Context) (*orderResolver, error) {
	order, err := rp.L("order").(*rp.OrderRepository).FindByID(ctx, r.m.Order_id)
	if err != nil {
		return nil, err
	}
	return &orderResolver{order}, nil
}

func (r *commentResolver) Content() string {
	return r.m.Content
}

func (r *commentResolver) CreatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Created_at)
	return graphql.Time{Time: res}, err
}

func (r *commentResolver) Reply() *string {
	return r.m.Reply
}

func (r *commentResolver) ReplyTime() (*graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, *r.m.Reply_time)
	return &graphql.Time{Time: res}, err
}

func (r *commentResolver) IsShow() bool {
	return r.m.Is_show
}

func (r *commentResolver) UserIp() *string {
	return r.m.User_ip
}

func (r *commentResolver) ShippingRank() *int32 {
	return r.m.Shipping_rank
}

func (r *commentResolver) ServiceRank() *int32 {
	return r.m.Service_rank
}

func (r *commentResolver) Images(ctx context.Context) *[]string {
	if r.m.Images == nil {
		return nil
	}
	ra := util.Substr(*r.m.Images)
	res := strings.Split(ra, ",")
	for x, y := range res {
		//res[x] = util.TrimQuotes(y)
		res[x] = CompleteUrl(ctx, y)
	}
	return &res
}

func (r *commentResolver) UpdatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Updated_at)
	return graphql.Time{Time: res}, err
}

func (r *commentResolver) IsAnonymous() bool {
	return r.m.Is_anonymous
}
