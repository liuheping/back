package context

import (
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/model"
)

func CurrentUser(ctx context.Context) *string {
    return ctx.Value("user_id").(*string)
}

func User(ctx context.Context) *model.User {
    return ctx.Value("user").(*model.User)
}

func IsAlly(ctx context.Context) bool {
    return ctx.Value("is_authorized").(bool) && User(ctx).Type == model.ALLY
}

func CheckAuth(ctx context.Context) {
    if !ctx.Value("is_authorized").(bool) {
        panic(CredentialsError)
    }
}
