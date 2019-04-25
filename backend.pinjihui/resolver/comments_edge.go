package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"pinjihui.com/backend.pinjihui/model"
)

type commentsEdgeResolver struct {
	cursor graphql.ID
	m      *model.Comment
}

func (r *commentsEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *commentsEdgeResolver) Node() *commentResolver {
	return &commentResolver{r.m}
}
