package routes

import (
	"github.com/MashuNakamura/todolist-backend/controllers"
	"github.com/MashuNakamura/todolist-backend/middleware"
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
	api.Get("/profile", controllers.GetProfile)              // Read One User
	api.Post("/logout", controllers.Logout)                  // Logout

	// Protected Route
	protected := api.Group("/", middleware.Protected)

	protected.Post("/change-password", controllers.ChangePassword) // Change Password

	// Task API Route
	protected.Post("/tasks", controllers.CreateTask)                     // Create
	protected.Get("/tasks", controllers.GetAllTasks)                     // Read All
	protected.Get("/tasks/:id", controllers.GetTaskByID)                 // Read One
	protected.Get("/tasks/user/:user_id", controllers.GetAllTasksByUser) // Read All by User
	protected.Put("/tasks/:id", controllers.UpdateTask)                  // Update
	protected.Delete("/tasks/:id", controllers.DeleteTask)               // Delete
	protected.Put("/tasks/status", controllers.UpdateBatchStatus)        // Update Batch Status

	// Category API Route
	protected.Post("/categories", controllers.CreateCategory)                   // Create
	protected.Get("/categories/user/:user_id", controllers.GetCategoriesByUser) // Read All
	protected.Put("/categories/:id", controllers.UpdateCategory)                // Update
	protected.Delete("/categories/:id", controllers.DeleteCategory)             // Delete
}
