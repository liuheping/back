package main

import (
    "log"
    "net/http"

    gcontext "pinjihui.com/pinjihui/context"
    h "pinjihui.com/pinjihui/handler"
    "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/resolver"
    "pinjihui.com/pinjihui/schema"

    "time"

    "github.com/graph-gophers/graphql-go"
    "github.com/nicksnyder/go-i18n/i18n"
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/loader"
    "pinjihui.com/pinjihui/service"
)

func main() {
    i18n.MustLoadTranslationFile("./i18n/zh-CN.all.json")
    config := gcontext.LoadConfig(".")

    db, err := gcontext.OpenDB(config)
    if err != nil {
        log.Fatalf("Unable to connect to db: %s \n", err)
    }
    ctx := context.Background()
    log := service.NewLogger(config)
    repository.RegisterAll(db, log)
    authService := service.NewAuthService(config, log)

    ctx = context.WithValue(ctx, "config", config)
    ctx = context.WithValue(ctx, "log", log)
    ctx = context.WithValue(ctx, "db", db)
    ctx = context.WithValue(ctx, "authService", authService)

    graphqlSchema := graphql.MustParseSchema(schema.GetRootSchema(), &resolver.Resolver{})

    http.Handle("/login", h.AddContext(ctx, h.Login()))
    http.Handle("/wx_notify", h.AddContext(ctx, h.WxNotify()))

    loggerHandler := &h.LoggerHandler{config.DebugMode}
    http.Handle("/query", h.AddContext(ctx, loggerHandler.Logging(h.Authenticate(&h.GraphQL{Schema: graphqlSchema, Loaders: loader.NewLoaderCollection()}))))

    http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "graphiql.html")
    }))
    http.Handle("/test_upload.html", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "test_upload.html")
    }))

    go doSchedule()
    log.Fatal(http.ListenAndServe(config.ListenPort, nil))
}

func doSchedule() {
    repository.L("spike").(*repository.SpikeRepository).UpdateProductPrice()
    repository.L("order").(*repository.OrderRepository).UpdateOrderStatusBySchedule()
    time.Sleep(3 * time.Second)
    doSchedule()
}
