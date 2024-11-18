package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"isit/isit/internal/service/models"
	"isit/isit/pkg/validator"
	"log/slog"
)

type Service interface{}

type ServiceHandler struct {
	repo *sqlx.DB
	app  *fiber.App
}

func NewServiceHandler(app *fiber.App, repo *sqlx.DB) *ServiceHandler {
	return &ServiceHandler{
		repo: repo,
		app:  app,
	}
}

func (h *ServiceHandler) InitRoute() {
	h.app.Get("/material", h.GetMaterials())
	h.app.Get("/product", h.GetProducts())
	h.app.Post("/material", h.AddMaterial())
	h.app.Post("/product", h.AddProduct())
	h.app.Patch("/material", h.UpdateMaterial())

	h.app.Group("/order")
	h.app.Post("/", h.AddOrder())
	h.app.Get("/", h.GetOrders())
	h.app.Patch("/status", h.UpdateStatus())

}

func (h *ServiceHandler) AddMaterial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.Material)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(err)
		}

		query := `
INSERT INTO materials.material
(name, price, description, address, count, reject_percent, sending_date)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
`
		var result int

		err = h.repo.GetContext(c.Context(), &result, query,
			params.Name, params.Price, params.Description, params.Address, params.Count, params.RejectPercent)
		if err != nil {
			slog.Error("failed add material", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"id": result,
		})
	}
}

func (h *ServiceHandler) AddProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.Product)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(err)
		}

		tx, err := h.repo.BeginTxx(c.Context(), nil)
		if err != nil {
			slog.Error("failed begin tx", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(err)
		}

		defer tx.Rollback()

		query := `
INSERT INTO products.product
(name, description, price, reject_percent)
VALUES ($1, $2, $3, $4)
RETURNING id;
`
		var result int

		err = tx.GetContext(c.Context(), &result, query,
			params.Name, params.Description, params.Price, params.RejectPercent)
		if err != nil {
			slog.Error("failed add material", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		query = `
INSERT INTO products.material
(product_id, material_id, count)
SELECT unnest($1::int), unnest($2::numeric)
`
		ids := make([]int, 0)
		counts := make([]float64, 0)
		for _, item := range params.Materials {
			ids = append(ids, item.Id)
			counts = append(counts, item.Count)
		}

		_, err = tx.ExecContext(c.Context(), query, ids, counts)
		if err != nil {
			slog.Error("failed add material", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		tx.Commit()

		return c.JSON(fiber.Map{
			"id": result,
		})
	}
}

func (h *ServiceHandler) AddOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.Order)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(err)
		}

		query := `
INSERT INTO orders.order
(product_id, count, status)
VALUES ($1, $2, $3)
`
		var result int

		err = h.repo.GetContext(c.Context(), &result, query,
			params.ProductId, params.Count, "created")
		if err != nil {
			slog.Error("failed add material", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"id": result,
		})
	}
}

func (h *ServiceHandler) UpdateMaterial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.Material)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(err)
		}

		query := `
UPDATE materials.material
SET  name = $1,
price = $2,
description = $3,
address = $4,
reject_percent = $5,
count = $6
WHERE id = $7 
`
		var result int

		_, err = h.repo.ExecContext(c.Context(), query,
			params.Name, params.Price, params.Description, params.Address,
			params.RejectPercent, params.Count, params.Count)
		if err != nil {
			slog.Error("failed update material", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"id": result,
		})
	}
}

func (h *ServiceHandler) UpdateStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.UpdateOrder)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(err)
		}

		query := `
UPDATE orders."order"
SET status = $2
WHERE id = $1
`
		var result int

		_, err = h.repo.ExecContext(c.Context(), query, params.Id, params.Status)
		if err != nil {
			slog.Error("failed update order", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"id": result,
		})
	}
}

func (h *ServiceHandler) GetProducts() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.UpdateOrder)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(err)
		}

		query := `
Select  p.id,
p.name,
p.description,
p.price,
p.reject_percent,
array(select m.material_id from products.material m where m.product_id = p.id) AS materials
    from products.product p
`
		result := make([]models.Orders, 0)

		err = h.repo.SelectContext(c.Context(), &result, query)
		if err != nil {
			slog.Error("failed update order", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		return c.JSON(result)
	}
}

func (h *ServiceHandler) GetMaterials() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.UpdateOrder)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(err)
		}

		query := `
Select  id,
name,
price,
description,
address,
reject_percent,
count
    from materials.material
`
		result := make([]models.Material, 0)

		err = h.repo.SelectContext(c.Context(), &result, query)
		if err != nil {
			slog.Error("failed update order", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		return c.JSON(result)
	}
}

func (h *ServiceHandler) GetOrders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.UpdateOrder)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(err)
		}

		query := `
Select  id,
product_id,
count,
status
    from orders.order
`
		result := make([]models.GetOrder, 0)

		err = h.repo.SelectContext(c.Context(), &result, query)
		if err != nil {
			slog.Error("failed get order", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		return c.JSON(result)
	}
}
