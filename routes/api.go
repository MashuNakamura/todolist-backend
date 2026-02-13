package routes

import (
	"github.com/MashuNakamura/todolist-backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// API Route
	api := app.Group("/api") // API Route

	// Health Check API Route
	api.Get("/health", controllers.HealthCheck) // Health Check

	// Task API Route
	api.Post("/tasks", controllers.CreateTask)                     // Create
	api.Get("/tasks", controllers.GetAllTasks)                     // Read All
	api.Get("/tasks/:id", controllers.GetTaskByID)                 // Read One
	api.Get("/tasks/user/:user_id", controllers.GetAllTasksByUser) // Read All by User
	api.Put("/tasks/:id", controllers.UpdateTask)                  // Update
	api.Delete("/tasks/:id", controllers.DeleteTask)               // Delete
	api.Put("/tasks/status", controllers.UpdateBatchStatus)        // Update Batch Status

	// User API Route
	api.Post("/register", controllers.Register) // Register
	api.Post("/login", controllers.Login)       // Login
}
