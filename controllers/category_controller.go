package controllers

import (
	"github.com/MashuNakamura/todolist-backend/config"
	"github.com/MashuNakamura/todolist-backend/models"
	"github.com/gofiber/fiber/v2"
)

// API Untuk Create Category
func CreateCategory(c *fiber.Ctx) error {
	var cat models.Category
	if err := c.BodyParser(&cat); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid input",
			Error:   400,
		})
	}

	if cat.Name == "" {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Name is required",
			Error:   400,
		})
	}

	if cat.Color == "" {
		cat.Color = "#000000"
	}

	if err := config.DB.Create(&cat).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to create category",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Category created successfully",
		Error:   200,
		Data:    cat,
	})
}

// API Untuk Get Category
func GetCategoriesByUser(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var cats []models.Category

	if err := config.DB.Where("user_id = ?", userID).Find(&cats).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to retrieve categories",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Categories retrieved successfully",
		Error:   200,
		Data:    cats,
	})
}

// API Untuk Delete Category
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var cat models.Category

	if err := config.DB.First(&cat, id).Error; err != nil {
		return c.Status(404).JSON(models.Ret{
			Success: false,
			Message: "Category not found",
			Error:   404,
		})
	}

	if err := config.DB.Delete(&cat).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to delete category",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Category deleted successfully",
		Error:   200,
	})
}

// API Untuk Update Category
func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var cat models.Category

	if err := config.DB.First(&cat, id).Error; err != nil {
		return c.Status(404).JSON(models.Ret{
			Success: false,
			Message: "Category not found",
			Error:   404,
		})
	}

	if err := c.BodyParser(&cat); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid input",
			Error:   400,
		})
	}

	if cat.Name == "" {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Name is required",
			Error:   400,
		})
	}

	if cat.Color == "" {
		cat.Color = "#000000"
	}

	if err := config.DB.Save(&cat).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to update category",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Category updated successfully",
		Error:   200,
		Data:    cat,
	})
}
