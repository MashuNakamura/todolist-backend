package controllers

import (
	"github.com/MashuNakamura/todolist-backend/config"
	"github.com/MashuNakamura/todolist-backend/models"
	"github.com/gofiber/fiber/v2"
)

type Ret struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   int    `json:"code"`
	Data    any    `json:"data"`
}

func CreateTask(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.JSON(Ret{
			Success: false,
			Message: "Failed to parse request body",
			Error:   400,
			Data:    nil,
		})
	}

	if err := config.DB.Create(&task).Error; err != nil {
		return c.JSON(Ret{
			Success: false,
			Message: "Failed to create task",
			Error:   500,
			Data:    nil,
		})
	}

	return c.JSON(Ret{
		Success: true,
		Message: "Task created successfully",
		Error:   200,
		Data:    task,
	})
}

func GetAllTasks(c *fiber.Ctx) error {
	var tasks []models.Task

	if err := config.DB.Find(&tasks).Error; err != nil {
		return c.JSON(Ret{
			Success: false,
			Message: "Failed to get tasks",
			Error:   500,
			Data:    nil,
		})
	}

	return c.JSON(Ret{
		Success: true,
		Message: "Tasks retrieved successfully",
		Error:   200,
		Data:    tasks,
	})
}
