package resolver

import (
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/model"
    rp "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/loader"
)

func (*Resolver) CreateComment(ctx context.Context, args struct{ NewComment *model.CommentInput }) (*commentResolver, error) {
    comment, err := rp.L("comment").(*rp.CommentRepository).Create(ctx, args.NewComment)
    if err != nil {
        return nil, err
    }
    product, err := loader.LoadProduct(ctx, comment.ProductID, comment.MerchantID)
    if err != nil {
        return nil, err
    }
    return &commentResolver{comment, &product.Product}, nil
}
