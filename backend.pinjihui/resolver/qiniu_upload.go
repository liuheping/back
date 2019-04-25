package resolver

import (
	"golang.org/x/net/context"
	gc "pinjihui.com/backend.pinjihui/context"
	"pinjihui.com/backend.pinjihui/service"
)

type QiniuUploadTokenResolver struct {
	token *service.QiniuUploadToken
}

func (q *QiniuUploadTokenResolver) Token() string {
	return q.token.Token
}

// func (q *QiniuUploadTokenResolver) Key() string {
// 	return q.token.Key
// }

func (r *Resolver) QiniuUploadToken(ctx context.Context, args struct {
	Module string
	Ext    string
}) *QiniuUploadTokenResolver {
	gc.CheckAuth(ctx)
	token := service.GetQiniuUploadToken(ctx, args.Module, args.Ext)
	return &QiniuUploadTokenResolver{token}
}
