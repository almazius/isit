package http

import (
	"github.com/gofiber/fiber/v2"
	"isit/internal/auth/service"
	"log/slog"
)

type AuthMW struct {
	authService *service.AuthService
}

func NewAuthMW(authService *service.AuthService) *AuthMW {
	return &AuthMW{authService: authService}
}

func (mw *AuthMW) AuthedMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionKey := c.Cookies("Authorization")
		userAgent := c.Get("User-Agent")
		fingerprint := c.Get("Fingerprint")

		if len(sessionKey) == 0 || len(fingerprint) == 0 || len(userAgent) == 0 {
			slog.Error("session key or fingerprint or user-agent is nil")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		err := mw.authService.ValidateSessionKey(c.UserContext(), userAgent, fingerprint, sessionKey)
		if err != nil {
			slog.Error("failed to validate session key", "error", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.Next()
	}
}
