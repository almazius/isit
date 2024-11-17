package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"isit/isit/config"
	"isit/isit/internal/auth/http"
	"isit/isit/internal/auth/repository/postgtres"
	redis2 "isit/isit/internal/auth/repository/redis"
	"isit/isit/internal/auth/service"
	"isit/isit/internal/role/route"
	"isit/isit/internal/role/usecase"
	http2 "isit/isit/internal/service/http"
	"isit/isit/pkg/logger"
	"isit/isit/pkg/postrges"
	"isit/isit/pkg/redis"
	"log/slog"
)

func main() {
	logger.Init()

	v, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed load config", "error", err)
		return
	}

	err = config.ParseConfig(v)
	if err != nil {
		slog.Error("failed parse config", "error", err)
		return
	}

	Init()
}

func Init() {
	db, err := postrges.InitPsqlDB(config.C())
	if err != nil {
		slog.Error("failed init postgres db", "error", err)
		panic(err)
	}

	//err = automigrations.InitRepository(context.TODO(), db)
	//if err != nil {
	//	slog.Error("failed init automigrations", "error", err)
	//	panic(err)
	//}

	rdb, err := redis.Init(config.C())
	if err != nil {
		slog.Error("failed init redis db", "error", err)
		panic(err)
	}

	authCache := redis2.NewAuthCache(rdb)
	authRepository := postgtres.NewAuthRepository(db)

	authService := service.NewAuthService(authRepository, authCache)
	roleService := usecase.NewRoleService()

	app := fiber.New()

	roleHandlers := route.NewRoleHandlers(app, roleService)
	authHandler := http.NewAuthHandler(authService)
	service := http2.NewServiceHandler(app, db)
	service.InitRoute()

	mw := http.NewAuthMW(authService)

	authGroup := app.Group("/auth")
	apiGroup := app.Group("/api")
	roleHandlers.InitRoleMap(mw)

	http.AuthRoute(authGroup, authHandler)
	http2.ServiceRoute(apiGroup, mw, nil)
	err = app.Listen(fmt.Sprintf(":%d", config.C().Port))
	if err != nil {
		slog.Error("failed listen port for start service", "error", err,
			slog.Int("port", config.C().Port))
		return
	}
}
