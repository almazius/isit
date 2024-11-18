package route

import (
	"github.com/gofiber/fiber/v2"
	"isit/internal/auth/http"
	"isit/internal/role/models"
	"isit/internal/role/usecase"
	"isit/pkg/validator"
	"log/slog"
)

type RoleHandlers struct {
	roleService *usecase.RoleService
	app         *fiber.App
}

func NewRoleHandlers(app *fiber.App, rs *usecase.RoleService) *RoleHandlers {
	return &RoleHandlers{
		app:         app,
		roleService: rs,
	}
}

func (h *RoleHandlers) InitRoleMap(mw *http.AuthMW) {
	route := h.app.Group("/permission")

	route.Get("/role", mw.AuthedMiddleware(), h.Validate(), h.GetRoles())
	route.Get("/permission", mw.AuthedMiddleware(), h.Validate(), h.GetPermissions())

	route.Post("/permission", mw.AuthedMiddleware(), h.Validate(), h.CreatePermission())
	route.Post("/role", mw.AuthedMiddleware(), h.Validate(), h.CreateRole())

	route.Patch("/grant_permission", mw.AuthedMiddleware(), h.Validate(), h.GrantPermission())
}

func (h *RoleHandlers) Validate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionKey := c.Cookies("Authorization")

		_ = sessionKey
		// get role info by session

		return c.Next()
	}
}

func (h *RoleHandlers) CreateRole() fiber.Handler {
	return func(c *fiber.Ctx) error {

		params := models.Role{}

		err := validator.ParseBody(c, &params)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = h.roleService.CreateRole(c.UserContext(), &params)
		if err != nil {
			slog.Error("failed create role", "error", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func (h *RoleHandlers) GetRoles() fiber.Handler {
	return func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"roles": []string{
				"professor",
				"student",
				"admin",
			},
		})
	}
}

func (h *RoleHandlers) GetPermissions() fiber.Handler {
	return func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"permissions": []models.Permission{
				models.Permission{
					ID:     0,
					Name:   "get permissions",
					Path:   "/permissions/permission",
					Method: "GET",
				},
				models.Permission{
					ID:     1,
					Name:   "create permissions",
					Path:   "/permissions/permission",
					Method: "POST",
				},
				models.Permission{
					ID:     2,
					Name:   "get roles",
					Path:   "/permissions/role",
					Method: "GET",
				},
				models.Permission{
					ID:     3,
					Name:   "create role",
					Path:   "/permissions/role",
					Method: "POST",
				},
				models.Permission{
					ID:     4,
					Name:   "grant permission",
					Path:   "/permissions/permission",
					Method: "PATCH",
				},
			},
		})
	}
}

func (h *RoleHandlers) CreatePermission() fiber.Handler {
	return func(c *fiber.Ctx) error {

		params := models.Permission{}

		err := validator.ParseBody(c, &params)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = h.roleService.CreatePermission(c.UserContext(), &params)
		if err != nil {
			slog.Error("failed create role", "error", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func (h *RoleHandlers) GrantPermission() fiber.Handler {
	return func(c *fiber.Ctx) error {

		params := models.GrantPermission{}

		err := validator.ParseBody(c, &params)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = h.roleService.GrantPermission(c.UserContext(), &params)
		if err != nil {
			slog.Error("failed create role", "error", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
