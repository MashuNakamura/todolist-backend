package routes

import (
	"github.com/MashuNakamura/todolist-backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Backend API Route
	api.Post("/tasks", controllers.CreateTask)
	api.Get("/tasks", controllers.GetAllTasks)
}
