package main

import (
	"CMS-Backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes() *fiber.App {
	router := fiber.New()
	router.Get("/login", handlers.Login)
	router.Get("/register", handlers.Register)
	router.Get("/logout", handlers.Logout)
	return router
}