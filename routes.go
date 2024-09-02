package main

import (
	"CMS-Backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/login", handlers.Login)
	router.Post("/register", handlers.Register)
	router.Get("/logout", handlers.Logout)
	router.Get("/validate-session", handlers.ValidateSession)
	return router
}