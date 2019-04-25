package main

import (
	"fmt"
	"log"
	"net/http"

	gcontext "pinjihui.com/backend.pinjihui/context"
	h "pinjihui.com/backend.pinjihui/handler"
	"pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/resolver"
	"pinjihui.com/backend.pinjihui/schema"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/loader"
	"pinjihui.com/backend.pinjihui/service"
)

func main() {
	config := gcontext.LoadConfig(".")

	db, err := gcontext.OpenDB(config)
	if err != nil {
		log.Fatalf("Unable to connect to db: %s \n", err)
	}
	ctx := context.Background()
	log := service.NewLogger(config)
	authService := service.NewAuthService(config, log)
	repository.RegisterAll(db, log)
	roleRepository := repository.NewRoleRepository(db, log)
	userRepository := repository.NewUserRepository(db, roleRepository, log)

	config2 := viper.New()
	config2.SetConfigName("Config")
	config2.AddConfigPath(".")
	err2 := config2.ReadInConfig()
	if err2 != nil {
		log.Fatalf("Fatal error context file: %s \n", err)
	}
	ctx = context.WithValue(ctx, "config2", config2)

	ctx = context.WithValue(ctx, "config", config)
	ctx = context.WithValue(ctx, "log", log)
	ctx = context.WithValue(ctx, "roleRepository", roleRepository)
	ctx = context.WithValue(ctx, "userRepository", userRepository)
	ctx = context.WithValue(ctx, "authService", authService)

	graphqlSchema := graphql.MustParseSchema(schema.GetRootSchema(), &resolver.Resolver{})

	http.Handle("/login", h.AddContext(ctx, h.Login()))

	loggerHandler := &h.LoggerHandler{config.DebugMode}
	http.Handle("/query", h.AddContext(ctx, loggerHandler.Logging(h.Authenticate(&h.GraphQL{Schema: graphqlSchema, Loaders: loader.NewLoaderCollection()}))))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "graphiql.html")
	}))

	fmt.Println("API server listening at", config.ListenPort)
	log.Fatal(http.ListenAndServe(config.ListenPort, nil))

}
