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

	// User API Route
	api.Post("/register", controllers.Register)              // Register
	api.Post("/login", controllers.Login)                    // Login
	api.Post("/forgot-password", controllers.ForgotPassword) // Forgot Password
	api.Post("/reset-password", controllers.ResetPassword)   // Reset Password
	api.Post("/change-password", controllers.ChangePassword) // Change Password

	// Task API Route
	api.Post("/tasks", controllers.CreateTask)                     // Create
	api.Get("/tasks", controllers.GetAllTasks)                     // Read All
	api.Get("/tasks/:id", controllers.GetTaskByID)                 // Read One
	api.Get("/tasks/user/:user_id", controllers.GetAllTasksByUser) // Read All by User
	api.Put("/tasks/:id", controllers.UpdateTask)                  // Update
	api.Delete("/tasks/:id", controllers.DeleteTask)               // Delete
	api.Put("/tasks/status", controllers.UpdateBatchStatus)        // Update Batch Status

	// Category API Route
	api.Post("/categories", controllers.CreateCategory)                   // Create
	api.Get("/categories/user/:user_id", controllers.GetCategoriesByUser) // Read All
	api.Put("/categories/:id", controllers.UpdateCategory)                // Update
	api.Delete("/categories/:id", controllers.DeleteCategory)             // Delete
}
