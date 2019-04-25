package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

//根据条件获取所有评论
func (r *Resolver) Comments(ctx context.Context, args struct {
	First  *int32
	Offset *int32
	Search *model.CommentSearchInput
	Sort   *model.CommentSortInput
}) (*commentsConnectionResolver, error) {
	//decodedIndex, _ := service.DecodeCursor(args.After)
	fetchSize := util.GetInt32(args.First, DefaultPageSize)
	list, err := rp.L("comment").(*rp.CommentRepository).Search(ctx, &fetchSize, args.Offset, args.Search, args.Sort)
	if err != nil {
		return nil, err
	}
	count, err := rp.L("comment").(*rp.CommentRepository).Count(ctx, args.Search)
	if err != nil {
		return nil, err
	}
	var from *string
	var to *string
	if len(list) > 0 {
		from = &(list[0].ID)
		to = &(list[len(list)-1].ID)
	}
	return &commentsConnectionResolver{m: list, totalCount: count, from: from, to: to, hasNext: util.If(len(list) == int(fetchSize), true, false).(bool)}, nil
}

//根据ID查找评论
func (r *Resolver) Comment(ctx context.Context, args struct {
	ID string
}) (*commentResolver, error) {
	comment, err := rp.L("comment").(*rp.CommentRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &commentResolver{comment}, nil
}

// 通过商品ID查找评论
func (r *Resolver) FindCommentsByProductID(ctx context.Context, args struct {
	ProductID string
}) (*[]*commentResolver, error) {
	comments, err := rp.L("comment").(*rp.CommentRepository).FindByProductID(ctx, args.ProductID)
	if err != nil {
		return nil, err
	}
	l := make([]*commentResolver, len(comments))
	for i := range l {
		l[i] = &commentResolver{(comments)[i]}
	}
	return &l, nil
}
