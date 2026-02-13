package controllers

import (
	"github.com/MashuNakamura/todolist-backend/config"
	"github.com/MashuNakamura/todolist-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

// API Create Task
func CreateTask(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Failed to parse request body",
			Error:   400,
			Data:    nil,
		})
	}

	if task.Title == "" {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Task title is required",
			Error:   400,
			Data:    nil,
		})
	}

	if task.Priority == "" {
		task.Priority = "Medium"
	}

	if task.Status == "" {
		task.Status = "todo"
	}

	if err := config.DB.Create(&task).Error; err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Failed to create task",
			Error:   500,
			Data:    nil,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Task created successfully",
		Error:   200,
		Data:    task,
	})
}

// API Get All Tasks
func GetAllTasks(c *fiber.Ctx) error {
	var tasks []models.Task

	if err := config.DB.Find(&tasks).Error; err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Failed to get tasks",
			Error:   500,
			Data:    nil,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Tasks retrieved successfully",
		Error:   200,
		Data:    tasks,
	})
}

// API GetAllTasks By Specify User
func GetAllTasksByUser(c *fiber.Ctx) error {
	var tasks []models.Task

	if err := config.DB.Where("user_id = ?", c.Params("user_id")).Find(&tasks).Error; err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Failed to get tasks",
			Error:   500,
			Data:    nil,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Tasks retrieved successfully",
		Error:   200,
		Data:    tasks,
	})
}

// API Get Task by ID
func GetTaskByID(c *fiber.Ctx) error {
	var task models.Task

	if err := config.DB.Where("id = ?", c.Params("id")).First(&task).Error; err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Failed to get task",
			Error:   404,
			Data:    nil,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Task retrieved successfully",
		Error:   200,
		Data:    task,
	})
}

// API Edit Task by ID
func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")

	var task models.Task
	if err := config.DB.Where("id = ?", id).First(&task).Error; err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Failed to get task",
			Error:   404,
			Data:    nil,
		})
	}

	var updateTask models.UpdateTask
	if err := c.BodyParser(&updateTask); err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Failed to parse request body",
			Error:   400,
			Data:    nil,
		})
	}

	if updateTask.Title != "" {
		task.Title = updateTask.Title
	}
	if updateTask.ShortDesc != "" {
		task.ShortDesc = updateTask.ShortDesc
	}
	if updateTask.LongDesc != "" {
		task.LongDesc = updateTask.LongDesc
	}
	if updateTask.Priority != "" {
		task.Priority = updateTask.Priority
	}
	if updateTask.Status != "" {
		validStatuses := map[string]bool{"todo": true, "ongoing": true, "done": true}
		if validStatuses[updateTask.Status] {
			task.Status = updateTask.Status
		} else {
			return c.Status(400).JSON(models.Ret{
				Success: false,
				Message: "Invalid status (todo, ongoing, done)",
				Error:   400,
			})
		}
	}
	if updateTask.Time != "" {
		task.Time = updateTask.Time
	}
	if updateTask.Date != "" {
		task.Date = updateTask.Date
	}
	if updateTask.Tags != nil {
		task.Tags = pq.StringArray(updateTask.Tags)
	}

	if err := config.DB.Save(&task).Error; err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Failed to update task",
			Error:   500,
			Data:    nil,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Task updated successfully",
		Error:   200,
		Data:    task,
	})
}

// API Delete Task by ID (One or Many)
func DeleteTask(c *fiber.Ctx) error {
	var req models.DeleteTask

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid data",
			Error:   400,
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "No IDs provided",
			Error:   400,
		})
	}

	if err := config.DB.Where("id IN ?", req.IDs).Delete(&models.Task{}).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to delete tasks",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Tasks deleted successfully",
		Error:   200,
	})
}

// API Untuk Update One or Many Status
func UpdateBatchStatus(c *fiber.Ctx) error {
	var req models.UpdateBatchStatus
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Invalid data",
			Error:   400,
		})
	}

	if len(req.IDs) == 0 {
		return c.JSON(models.Ret{
			Success: false,
			Message: "No IDs provided",
			Error:   400,
		})
	}

	if err := config.DB.Model(&models.Task{}).Where("id IN ?", req.IDs).Update("status", req.Status).Error; err != nil {
		return c.JSON(models.Ret{
			Success: false,
			Message: "Failed to update tasks",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Tasks updated successfully",
		Error:   200,
	})
}
