package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ParseBody(c *fiber.Ctx, out interface{}) error {
	err := c.BodyParser(out)
	if err != nil {
		return err
	}

	err = validate.StructCtx(c.Context(), out)
	return err
}
