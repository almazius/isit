package http

import (
	"encoding/csv"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"isit/internal/auth/http"
	"isit/internal/service/models"
	"isit/pkg/validator"
	"log/slog"
	"os"
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

func (h *ServiceHandler) InitRoute(mw *http.AuthMW) {
	api := h.app.Group("/api")
	api.Get("/material", h.GetMaterials())
	api.Get("/product", h.GetProducts())
	api.Get("/material/csv", h.GetMaterialsCSV())
	api.Get("/product/csv", h.GetProductsCSV())
	api.Post("/material", h.AddMaterial())
	api.Post("/product", h.AddProduct())
	api.Patch("/material", h.UpdateMaterial())

	order := api.Group("/order")
	order.Post("/", h.AddOrder())
	order.Get("/", h.GetOrders())
	order.Get("/csv", h.GetOrdersCSV())

	order.Patch("/status", h.UpdateStatus())

}

// @Summary 	Добавление материала
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Param       body   body      models.Material  true "52"
// @Success	200													"http.StatusOK"
// @Router /api/material [post]
func (h *ServiceHandler) AddMaterial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.Material)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(map[string]string{
				"error": err.Error(),
			})
		}

		query := `
INSERT INTO materials.material
(name, price, description, address, count, reject_percent)
VALUES ($1, $2, $3, $4, $5, $6)
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

// @Summary 	Добавление продукта
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Param       body   body      models.Product  true "52"
// @Success	200													"http.StatusOK"
// @Router /api/product [post]
func (h *ServiceHandler) AddProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.Product)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(map[string]string{
				"error": err.Error(),
			})
		}

		query := `
INSERT INTO products.product
(name, description, price, reject_percent)
VALUES ($1, $2, $3, $4)
RETURNING id;
`
		var result int

		err = h.repo.GetContext(c.Context(), &result, query,
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
SELECT $3, unnest($1::int[]), unnest($2::numeric[])
`
		ids := make([]int, 0)
		counts := make([]float64, 0)
		for _, item := range params.Materials {
			ids = append(ids, item.Id)
			counts = append(counts, item.Count)
		}

		_, err = h.repo.ExecContext(c.Context(), query, pq.Array(ids), pq.Array(counts), result)
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

// @Summary 	Добавление заявки
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Param       body   body      models.Order  true "52"
// @Success	200													"http.StatusOK"
// @Router /api/order/ [post]
func (h *ServiceHandler) AddOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.Order)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(map[string]string{
				"error": err.Error(),
			})
		}

		query := `
INSERT INTO orders.order
(product_id, count, status)
VALUES ($1, $2, $3)
RETURNING id
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

// @Summary 	Обновелние материала
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Param       body   body      models.Material  true "52"
// @Success	200													"http.StatusOK"
// @Router /api/material [patch]
func (h *ServiceHandler) UpdateMaterial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.Material)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(map[string]string{
				"error": err.Error(),
			})
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
			params.RejectPercent, params.Count, params.Id)
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

// @Summary 	Обновелние материала
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Param       body   body      models.UpdateOrder  true "52"
// @Success	200													"http.StatusOK"
// @Router /api/order/status [patch]
func (h *ServiceHandler) UpdateStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(models.UpdateOrder)

		err := validator.ParseBody(c, params)
		if err != nil {
			slog.Error("failed parse params", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(map[string]string{
				"error": err.Error(),
			})
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

// @Summary 	Ручка
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Success	200			"http.StatusOK"
// @Router /api/product [get]
func (h *ServiceHandler) GetProducts() fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := `
Select  p.id,
p.name,
p.description,
p.price,
p.reject_percent
    from products.product p
`
		result := make([]models.Product, 0)

		err := h.repo.SelectContext(c.Context(), &result, query)
		if err != nil {
			slog.Error("failed update order", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		for i := range result {
			result[i].Materials = make([]models.MaterialSmallInfo, 0)
			query = `select m.material_id AS id, m.count AS count from products.material m where m.product_id = $1`
			err = h.repo.SelectContext(c.Context(), &result[i].Materials, query, result[i].Id)
			if err != nil {
				slog.Error("failed get orders", "error", err)
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": "internal server error",
					"error":   err.Error(),
				})
			}
		}

		return c.JSON(result)
	}
}

// @Summary 	Ручка
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Success	200			"http.StatusOK"
// @Router /api/product/csv [get]
func (h *ServiceHandler) GetProductsCSV() fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := `
Select  p.id,
p.name,
p.description,
p.price,
p.reject_percent
    from products.product p
`
		result := make([]models.Product, 0)

		err := h.repo.SelectContext(c.Context(), &result, query)
		if err != nil {
			slog.Error("failed update order", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		for i := range result {
			result[i].Materials = make([]models.MaterialSmallInfo, 0)
			query = `select m.material_id AS id, m.count AS count from products.material m where m.product_id = $1`
			err = h.repo.SelectContext(c.Context(), &result[i].Materials, query, result[i].Id)
			if err != nil {
				slog.Error("failed get orders", "error", err)
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": "internal server error",
					"error":   err.Error(),
				})
			}
		}

		fileName := uuid.New().String()
		file, _ := os.Create(fileName)
		defer os.Remove(fileName)
		cFile := csv.NewWriter(file)
		cFile.Write([]string{"ID продукта", "Название", "Цена", "Описание", "Процент брака"})

		for _, el := range result {
			cFile.Write([]string{fmt.Sprintf("%v", el.Id), fmt.Sprintf("%v", el.Name), fmt.Sprintf("%v", el.Price),
				fmt.Sprintf("%v", el.Description), fmt.Sprintf("%v", el.RejectPercent),
			})
		}
		cFile.Flush()

		return c.Download(fileName, "products.csv")
	}
}

// @Summary 	Ручка
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Success	200			"http.StatusOK"
// @Router /api/material [get]
func (h *ServiceHandler) GetMaterials() fiber.Handler {
	return func(c *fiber.Ctx) error {

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

		err := h.repo.SelectContext(c.Context(), &result, query)
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

// @Summary 	Ручка
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Success	200			"http.StatusOK"
// @Router /api/material/csv [get]
func (h *ServiceHandler) GetMaterialsCSV() fiber.Handler {
	return func(c *fiber.Ctx) error {

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

		err := h.repo.SelectContext(c.Context(), &result, query)
		if err != nil {
			slog.Error("failed update order", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}

		fileName := uuid.New().String()
		file, _ := os.Create(fileName)
		defer os.Remove(fileName)
		cFile := csv.NewWriter(file)
		cFile.Write([]string{"ID материала", "Название", "Цена", "Описание", "Адрес", "Процент брака", "Количество"})

		for _, el := range result {
			cFile.Write([]string{fmt.Sprintf("%v", el.Id), fmt.Sprintf("%v", el.Name), fmt.Sprintf("%v", el.Price),
				fmt.Sprintf("%v", el.Description), fmt.Sprintf("%v", el.Address), fmt.Sprintf("%v", el.RejectPercent),
				fmt.Sprintf("%v", el.Count)})
		}
		cFile.Flush()

		return c.Download(fileName, "materials.csv")
	}
}

// @Summary 	Ручка
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Success	200			"http.StatusOK"
// @Router /api/order [get]
func (h *ServiceHandler) GetOrders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := `
Select  id,
product_id,
count,
status
    from orders.order
`
		result := make([]models.GetOrder, 0)

		err := h.repo.SelectContext(c.Context(), &result, query)
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

// @Summary 	Ручка
// @Security	AuthToken
// @Tags 		File
// @Accept		json
// @Produce		json
// @Success	200			"http.StatusOK"
// @Router /api/order/csv [get]
func (h *ServiceHandler) GetOrdersCSV() fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := `
Select  id,
product_id,
count,
status
    from orders.order
`
		result := make([]models.GetOrder, 0)

		err := h.repo.SelectContext(c.Context(), &result, query)
		if err != nil {
			slog.Error("failed get order", "error", err)
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "internal server error",
				"error":   err.Error(),
			})
		}
		fileName := uuid.New().String()
		file, _ := os.Create(fileName)
		defer os.Remove(fileName)
		cFile := csv.NewWriter(file)
		cFile.Write([]string{"ID продукта", "Количество", "Статус заказа"})

		for _, el := range result {
			cFile.Write([]string{fmt.Sprintf("%v", el.ProductId), fmt.Sprintf("%v", el.Count), fmt.Sprintf("%v", el.Status)})
		}
		cFile.Flush()

		return c.Download(fileName, "orders.csv")
	}
}
