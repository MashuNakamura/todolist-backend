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
	api.Post("/register", controllers.Register)                 // Register
	api.Post("/login", controllers.Login)                       // Login
	api.Post("/forgot-password", controllers.ForgotPassword)    // Forgot Password
	api.Post("/reset-password", controllers.ResetPassword)      // Reset Password
	api.Get("/auth/google/login", middleware.GoogleLogin)       // Redirect ke Google
	api.Get("/auth/google/callback", middleware.GoogleCallback) // Callback dari Google

	// Protected Route
	protected := api.Group("/", middleware.Protected)
	protected.Post("/logout", controllers.Logout)                  // Logout
	protected.Post("/update-profile", controllers.UpdateProfile)   // Update Profile
	protected.Get("/profile", controllers.GetProfile)              // Read One User
	protected.Post("/change-password", controllers.ChangePassword) // Change Password

	// Task API Route
	protected.Post("/tasks", controllers.CreateTask)              // Create
	protected.Get("/tasks", controllers.GetAllTasks)              // Read All
	protected.Get("/tasks/:id", controllers.GetTaskByID)          // Read One
	protected.Put("/tasks/:id", controllers.UpdateTask)           // Update
	protected.Delete("/tasks", controllers.DeleteTask)            // Delete Batch Task
	protected.Put("/tasks/status", controllers.UpdateBatchStatus) // Update Batch Status

	// Category API Route
	protected.Post("/categories", controllers.CreateCategory)       // Create
	protected.Get("/categories", controllers.GetCategoriesByUser)   // Read All
	protected.Put("/categories/:id", controllers.UpdateCategory)    // Update
	protected.Delete("/categories/:id", controllers.DeleteCategory) // Delete
}
