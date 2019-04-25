package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	rp "pinjihui.com/backend.pinjihui/repository"
)

func (r *Resolver) CommentIsShow(ctx context.Context, args struct {
	CommentID graphql.ID
}) (bool, error) {
	_, err := rp.L("comment").(*rp.CommentRepository).Visible(ctx, string(args.CommentID))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Resolver) CommentReply(ctx context.Context, args struct {
	ID      string
	Content string
}) (*commentResolver, error) {
	comment, err := rp.L("comment").(*rp.CommentRepository).Reply(ctx, args.ID, args.Content)
	if err != nil {
		return nil, err
	}
	return &commentResolver{comment}, nil
}
