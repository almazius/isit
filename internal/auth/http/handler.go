package http

import (
	"github.com/gofiber/fiber/v2"
	"isit/internal/auth/models"
	"isit/internal/auth/service"
	"isit/pkg/validator"
	"log/slog"
)

type AuthHandel struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandel {
	return &AuthHandel{authService: authService}
}

// @Summary 	Регистрация
// @Security	AuthToken
// @Tags 		Авторизация
// @Accept		json
// @Produce		json
// @Param 		body 	body 	models.RegisterParams true "1"
// @Success	200			"http.StatusOK"
// @Router /auth/register [post]
func (h *AuthHandel) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params models.RegisterParams
		err := validator.ParseBody(c, &params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = h.authService.CreateUser(c.UserContext(), &params)
		if err != nil {
			slog.Error("failed create user", "error", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// @Summary 	Авторизация
// @Security	AuthToken
// @Tags 		Авторизация
// @Accept		json
// @Produce		json
// @Param 		body 	body 	models.AuthParams true "1"
// @Success	200			"http.StatusOK"
// @Router /auth/login [post]
func (h *AuthHandel) Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params models.AuthParams

		err := validator.ParseBody(c, &params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		userAgent := c.Get("User-Agent")
		fingerPrint := c.Get("Fingerprint")
		//if userAgent == "" || fingerPrint == "" {
		//	slog.Error("user-agent of fingerprint is null")
		//	return c.SendStatus(fiber.StatusBadRequest)
		//}

		sessionKey, err := h.authService.SingIn(c.UserContext(), &params, userAgent, fingerPrint)
		if err != nil {
			slog.Error("failed to authenticate user", "error", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if len(sessionKey) == 0 {
			slog.Error("failed login")
			return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"Succeed": "false", "Reason": "incorrect data"})
		}

		c.Cookie(
			&fiber.Cookie{
				Name:     "Authorization",
				Value:    sessionKey,
				Path:     "/",
				HTTPOnly: true,
			},
		)

		return c.JSON(map[string]string{"Succeed": "true"})
	}
}
