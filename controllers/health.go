package controllers

import "github.com/gofiber/fiber/v2"

// API Health Check
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Backend is running!",
	})
}
