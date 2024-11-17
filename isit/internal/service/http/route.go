package http

import (
	"github.com/gofiber/fiber/v2"
	"isit/isit/internal/auth/http"
)

func ServiceRoute(app fiber.Router, mw *http.AuthMW, handler interface{}) {
	app.Get("/test", mw.AuthedMiddleware(), func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.SendString("Hello world!")
		}
	}())

	app.Post("/add-material")
}
