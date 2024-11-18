package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"isit/config"
	"isit/docs"
	"isit/internal/auth/http"
	"isit/internal/auth/repository/postgtres"
	redis2 "isit/internal/auth/repository/redis"
	"isit/internal/auth/service"
	"isit/internal/role/route"
	"isit/internal/role/usecase"
	http2 "isit/internal/service/http"
	"isit/pkg/logger"
	"isit/pkg/postrges"
	"isit/pkg/redis"
	"log/slog"
)

// swag init -g cmd/api/main.go

// @title Бэкенд Сервиса
// @version 0.0.02
// @contact.name Docs developer
// @contact.url https://t.me/sigy922

// ---@host 127.0.0.1:8080
// ---@BasePath  /api

// @securityDefinitions.apikey AuthToken
// @in header
// @name sessionid
// @description autorization token from auth_service

// ---@securityDefinitions.apikey AuthToken
// ---@in cookie
// ---@name Authorization
// ---@description autorization token from auth_service
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

	_, err = db.Exec(initMigrations)
	if err != nil {
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

	// Конфигурация сваггера
	docs.SwaggerInfo.Host = "127.0.0.1:8080"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	roleHandlers := route.NewRoleHandlers(app, roleService)
	authHandler := http.NewAuthHandler(authService)
	service := http2.NewServiceHandler(app, db)

	mw := http.NewAuthMW(authService)

	service.InitRoute(mw)
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

var initMigrations = `
create schema if not exists users;
create schema if not exists materials;
create schema if not exists orders;
create schema if not exists products;
create schema if not exists users;
create schema if not exists cls;


create table if not exists users."user"
(
    id            bigserial
        primary key,
    login         text                                      not null,
    password_hash text                                      not null,
    name          text                     default ''::text not null,
    surname       text                     default ''::text not null,
    created_at    timestamp with time zone default now()    not null,
    deleted_at    timestamp with time zone,
    is_baned      boolean                  default false
);

alter table users."user"
    owner to "user";

create unique index if not exists user_login_idx
    on users."user" (login);

create table if not exists cls.role
(
    id   bigserial
        primary key,
    name text not null
);

alter table cls.role
    owner to "user";

create table if not exists users.user_role
(
    id      bigserial
        primary key,
    user_id bigserial
        references users."user"
            on delete cascade,
    role_id bigserial
        references cls.role
            on delete cascade
);

alter table users.user_role
    owner to "user";

create unique index if not exists user_role_idx
    on users.user_role (user_id, role_id);

create table if not exists materials.material
(
    id             bigserial
        primary key,
    name           text                                   not null,
    price          numeric(10, 2)                         not null,
    description    text,
    created_at     timestamp with time zone default now() not null,
    deleted_at     timestamp with time zone,
    address        text                                   not null,
    reject_percent numeric(10, 2)                         not null,
    sending_date   timestamp with time zone,
    count          numeric                  default 0     not null
);

alter table materials.material
    owner to "user";

create table if not exists products.product
(
    id             bigserial
        primary key,
    name           text           not null,
    description    text,
    price          numeric(10, 2) not null,
    reject_percent numeric(10, 2) not null
);

alter table products.product
    owner to "user";

create table if not exists products.material
(
    id          bigserial
        primary key,
    product_id  bigserial
        references products.product
            on delete cascade,
    material_id bigserial
        references materials.material
            on delete cascade,
    count       numeric default 0 not null
);

alter table products.material
    owner to "user";

create table if not exists orders."order"
(
    id         bigserial
        primary key,
    product_id bigint                                 not null
        references products.product,
    count      numeric                                not null,
    created_at timestamp with time zone default now() not null,
    status     text
);

alter table orders."order"
    owner to "user";


`
