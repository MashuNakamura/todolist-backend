package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Inisialisasi Fiber
	app := fiber.New()

	// 2. Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// 3. Test Route
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Backend is running! 🚀",
		})
	})

	// 4. Jalankan Server
	log.Fatal(app.Listen(":8080"))
}
