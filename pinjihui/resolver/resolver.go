package resolver

import (
    "strings"
    "golang.org/x/net/context"
    gc "pinjihui.com/pinjihui/context"
)

type Resolver struct{}

func completeUrl(ctx context.Context, url string) string {
    if !strings.HasPrefix(url, "http") {
        url = ctx.Value("config").(*gc.Config).QiniuCndHost + url
        return url
    }
    return url
}

func getThumbnail(url string) string {
    return url + "?imageView2/0/w/200/h/200"
}