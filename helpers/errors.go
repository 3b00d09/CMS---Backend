package helpers

import (
	"github.com/gofiber/fiber/v3"
)

func ServerError(c fiber.Ctx, err error) error{
	return c.JSON(fiber.Map{
		"message": "Internal server error.",
		"detail":err.Error(),
	})
}

func SessionError(c fiber.Ctx) error{
	return c.JSON(fiber.Map{
		"message": "Invalid session.",
	})
}

func FormError(c fiber.Ctx) error{
	return c.JSON(fiber.Map{
		"message": "Missing form fields or query parameters.",
	})
}