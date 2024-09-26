package main

import (
	"CMS-Backend/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func SetupRoutes() *fiber.App {
	router := fiber.New()
	router.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Content-Length", "Accept-Language", "Accept-Encoding", "Connection", "Access-Control-Allow-Origin"},
		AllowOrigins:     []string{"https://cms-frontend-angular-gamma.vercel.app", "http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
	}))

	// auth routes
	router.Post("/login", handlers.Login)
	router.Post("/register", handlers.Register)
	router.Get("/logout", handlers.Logout)
	router.Get("/validate-session", handlers.ValidateSession)

	router.Post("/create-project", handlers.HandleCreateProject)
	router.Get("projects", handlers.HandleGetProjects)
	router.Get("/users", handlers.HandleSearchUsers)
	router.Post("/create-page", handlers.HandleCreatePage)
	router.Put("/update-page", handlers.HandleUpdatePage)
	router.Get("/page", handlers.HandleGetPage)
	
	// for total projects and pages created by user
	router.Get("/stats-page", handlers.HandleGetStatsPage)
	router.Get("/last-modified", handlers.HandleGetLastModified)
	router.Delete("/delete-page", handlers.HandleDeletePage)
	router.Delete("/delete-project", handlers.HandleDeleteProject)
	return router
}