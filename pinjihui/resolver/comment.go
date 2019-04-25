package resolver

import (
    "time"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/loader"
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/util"
    "strings"
    gc "pinjihui.com/pinjihui/context"
    "pinjihui.com/pinjihui/repository"
)

type commentResolver struct {
    m       *model.Comment
    product *model.Product
}

type CommentUserResolver struct {
    u *model.User
}

func (u *CommentUserResolver) Name() string {
    return util.GetString(u.u.Name, "***")
}

func (u CommentUserResolver) AvatarUrl() *string {
    return nil
}

func (r *commentResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *commentResolver) User(ctx context.Context) (*CommentUserResolver, error) {
    user, err := loader.LoadUser(ctx, r.m.UserID)
    if err != nil {
        return nil, err
    }
    return &CommentUserResolver{user}, nil
}

func (r *commentResolver) Rank() int32 {
    return r.m.Rank
}

func (r *commentResolver) Content() string {
    return r.m.Content
}

func (r *commentResolver) CreatedAt() (graphql.Time, error) {
    res, err := time.Parse(time.RFC3339, r.m.CreatedAt)
    return graphql.Time{Time: res}, err
}

func (r *commentResolver) Reply() *string {
    return r.m.Reply
}

func (r *commentResolver) Images(ctx context.Context) []string {
    urls := r.m.GetImageArr()
    for i, url := range urls {
        //生成七牛外链
        if !strings.HasPrefix(url, "http") {
            urls[i] = ctx.Value("config").(*gc.Config).QiniuCndHost + url
        }
    }
    return urls
}

func (r *commentResolver) ProductSpec() ([]*attributeResolver, error) {
    p, err := repository.L("order").(*repository.OrderRepository).FindOrderProduct(r.m.OrderID, r.m.ProductID)
    if err != nil {
        return nil, err
    }
    spec := model.Spec{p.Spec1Name, p.Spec2Name, p.Spec1, p.Spec2}
    return resolverSpec(&spec), nil
}
