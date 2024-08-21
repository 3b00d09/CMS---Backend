package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const port string = ":3000"

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})
	app.Listen(port)
	fmt.Printf("Server Running on http://localhost%s\n", port)
}