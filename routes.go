package main

import (
	"CMS-Backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes() *fiber.App {
	router := fiber.New()
	router.Use(cors.New(cors.Config{
	AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
	AllowOrigins:     "http://localhost:4200",
	AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	AllowCredentials: true,
	}))
	router.Post("/login", handlers.Login)
	router.Post("/register", handlers.Register)
	router.Get("/logout", handlers.Logout)
	router.Get("/validate-session", handlers.ValidateSession)
	return router
}