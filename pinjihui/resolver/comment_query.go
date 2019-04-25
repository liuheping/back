package resolver

import (
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/service"
    "pinjihui.com/pinjihui/util"
    "pinjihui.com/pinjihui/loader"
    "golang.org/x/net/context"
    rp "pinjihui.com/pinjihui/repository"
)

func (r *Resolver) CommentList(ctx context.Context, args struct {
    First      *int32
    After      *string
    ProductID  graphql.ID
    MerchantID graphql.ID
    Rank       *int32
}) (*commentsConnectionResolver, error) {
    decodedIndex, _ := service.DecodeCursor(args.After)
    fetchSize := util.GetInt32(args.First, DefaultPageSize)
    productID := string(args.ProductID)
    merchantID := string(args.MerchantID)
    list, err := rp.L("comment").(*rp.CommentRepository).List(fetchSize, decodedIndex, productID, merchantID, args.Rank)
    if err != nil {
        return nil, err
    }
    count, err := rp.L("comment").(*rp.CommentRepository).Count(productID, merchantID, args.Rank)
    if err != nil {
        return nil, err
    }
    var from *string
    var to *string
    if len(list) > 0 {
        from = &(list[0].ID)
        to = &(list[len(list)-1].ID)
    }
    product, err := loader.LoadProduct(ctx, productID, merchantID)
    if err != nil {
        return nil, err
    }
    return &commentsConnectionResolver{list, &product.Product, Connection{totalCount: count, from: from, to: to, hasNext: len(list) == int(fetchSize)}}, nil
}

func (r *Resolver) Comment(args struct{ ID string }) (*commentResolver, error) {
    comment, err := rp.L("comment").(*rp.CommentRepository).FindByID(args.ID)
    if err != nil {
        return nil, err
    }
    p, err := rp.L("product").(*rp.ProductRepository).FindByID(comment.ProductID)
    if err != nil {
        return nil, err
    }
    return &commentResolver{comment, p}, nil
}
