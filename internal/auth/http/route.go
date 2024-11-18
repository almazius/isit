package http

import (
	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app fiber.Router, handler *AuthHandel) {
	app.Post("/login", handler.Auth())
	app.Post("/register", handler.Register())
}
