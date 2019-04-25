package service

import (
    "encoding/base64"
    "fmt"
    "github.com/graph-gophers/graphql-go"
    "strings"
    "github.com/qiniu/api.v7/auth/qbox"
    "github.com/qiniu/api.v7/storage"
    "golang.org/x/net/context"
    gc "pinjihui.com/pinjihui/context"
    "github.com/rs/xid"
    "time"
)

func EncodeCursor(i *string) graphql.ID {
    if i == nil {
        return graphql.ID("")
    }
    return graphql.ID(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%s", *i))))
}

func DecodeCursor(after *string) (*string, error) {
    var decodedValue *string
    if after != nil {
        b, err := base64.StdEncoding.DecodeString(*after)
        if err != nil {
            return nil, err
        }
        i := strings.TrimPrefix(string(b), "cursor")
        decodedValue = &i
    }
    return decodedValue, nil
}

type QiniuUploadToken struct {
    Token string
    Key   string
}

func GetQiniuUploadToken(ctx context.Context, module, ext string) *QiniuUploadToken {
    // 简单上传凭证
    bucket := ctx.Value("config").(*gc.Config).QiniuBucket
    accessKey := ctx.Value("config").(*gc.Config).QiniuAccessKey
    secretKey := ctx.Value("config").(*gc.Config).QiniuSecretKey

    key := fmt.Sprintf("%s/%s/%s.%s", module, time.Now().Format("2006-01-02"), xid.New().String(), ext)
    putPolicy := storage.PutPolicy{
        Scope:      bucket,
        FsizeLimit: 1024 * 1024,
        SaveKey:    key,
        MimeLimit:  "image/jpeg;image/png",
    }
    mac := qbox.NewMac(accessKey, secretKey)
    upToken := putPolicy.UploadToken(mac)
    return &QiniuUploadToken{upToken, key}
}
