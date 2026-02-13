package controllers

import (
	"os"
	"time"

	"github.com/MashuNakamura/todolist-backend/config"
	"github.com/MashuNakamura/todolist-backend/helper"
	"github.com/MashuNakamura/todolist-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Failed to parse user",
			Error:   400,
		})
	}

	if input.Email == "" || input.Password == "" || input.Name == "" {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Name, Email, and Password are required",
			Error:   400,
		})
	}

	if !helper.IsValidEmail(input.Email) {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid email format",
			Error:   400,
		})
	}

	if !helper.IsStrongPassword(input.Password) {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character",
			Error:   400,
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to generate password hash",
			Error:   500,
		})
	}

	newUser := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Email already exists",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "User created successfully",
		Error:   200,
		Data:    newUser,
	})
}

func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid input",
			Error:   400,
		})
	}

	if input.Email == "" || input.Password == "" {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Email and password are required",
			Error:   400,
		})
	}

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(models.Ret{
			Success: false,
			Message: "Invalid email or password",
			Error:   401,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(models.Ret{
			Success: false,
			Message: "Invalid password",
			Error:   401,
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to generate token",
			Error:   500,
		})
	}

	user.Password = ""

	return c.JSON(models.Ret{
		Success: true,
		Message: "Login successfully",
		Error:   200,
		Data: fiber.Map{
			"user":    user,
			"token":   token,
			"expires": time.Now().Add(time.Hour * 24).Unix(),
		},
	})
}
