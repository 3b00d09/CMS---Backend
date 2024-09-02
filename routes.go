package main

import (
	"CMS-Backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes() *fiber.App {
	router := fiber.New()
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})
	router.Post("/login", handlers.Login)
	router.Post("/register", handlers.Register)
	router.Get("/logout", handlers.Logout)
	return router
}