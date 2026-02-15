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

	newUser.Password = ""

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

// API Untuk Update Profile
func UpdateProfile(c *fiber.Ctx) error {
	val := c.Locals("user_id")
	userID, ok := val.(uint)
	if !ok {
		return c.Status(401).JSON(models.Ret{
			Success: false,
			Message: "Unauthorized: Invalid User Session",
			Error:   401,
		})
	}

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid input",
			Error:   400,
		})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(models.Ret{
			Success: false,
			Message: "User not found",
			Error:   404,
		})
	}

	if user.Name == input.Name {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "No changes to update",
			Error:   400,
		})
	}

	user.Name = input.Name

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to update profile",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Profile updated successfully",
		Error:   200,
	})
}

// API Untuk Forgot Password
func ForgotPassword(c *fiber.Ctx) error {
	var input models.ForgotPassword
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid input",
			Error:   400,
		})
	}

	if input.Email == "" {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Email is required",
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

	var user_cp models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user_cp).Error; err != nil {
		return c.Status(404).JSON(models.Ret{
			Success: false,
			Message: "User not found",
			Error:   404,
		})
	}

	otp := helper.GenerateOTP()
	user_cp.OTP = otp
	user_cp.OTPExpiry = time.Now().Add(5 * time.Minute).Unix()

	if err := config.DB.Save(&user_cp).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to generate OTP",
			Error:   500,
		})
	}

	// Send Email
	emailBody := "Your OTP code is: " + otp
	if err := helper.SendEmail(user_cp.Email, "Reset Password OTP", emailBody); err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to send email",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "OTP has been sent to your email",
		Error:   200,
	})
}

// API Untuk Reset Password
func ResetPassword(c *fiber.Ctx) error {
	var input models.ResetPassword
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid input",
			Error:   400,
		})
	}

	var user models.User
	if err := config.DB.Where("email = ? AND otp = ?", input.Email, input.OTP).First(&user).Error; err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid Email or OTP",
			Error:   400,
		})
	}

	if time.Now().Unix() > user.OTPExpiry {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "OTP has expired",
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

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err == nil {
		return c.Status(400).JSON(models.Ret{Success: false, Message: "New password cannot be the same as old password", Error: 400})
	}

	hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to generate password hash",
			Error:   500,
		})
	}

	user.Password = string(hashNewPassword)
	user.OTP = ""
	user.OTPExpiry = 0

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to reset password",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Password reset successfully",
		Error:   200,
	})
}

// API Untuk Change Password
func ChangePassword(c *fiber.Ctx) error {
	val := c.Locals("user_id")
	userID, ok := val.(uint)
	if !ok {
		return c.Status(401).JSON(models.Ret{
			Success: false,
			Message: "Unauthorized: Invalid User Session",
			Error:   401,
		})
	}

	var input models.ChangePassword
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Invalid input",
			Error:   400,
		})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(models.Ret{
			Success: false,
			Message: "User not found",
			Error:   404,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Incorrect old password",
			Error:   400,
		})
	}

	if !helper.IsStrongPassword(input.NewPassword) {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character",
			Error:   400,
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.NewPassword))
	if err == nil {
		return c.Status(400).JSON(models.Ret{
			Success: false,
			Message: "New password cannot be the same as old password",
			Error:   400,
		})
	}

	hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 10)
	if err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to generate password hash",
			Error:   500,
		})
	}

	user.Password = string(hashNewPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(models.Ret{
			Success: false,
			Message: "Failed to change password",
			Error:   500,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "Password changed successfully",
		Error:   200,
	})
}

// API Untuk Get Profile
func GetProfile(c *fiber.Ctx) error {
	val := c.Locals("user_id")
	userID, ok := val.(uint)
	if !ok {
		return c.Status(401).JSON(models.Ret{
			Success: false,
			Message: "Unauthorized: Invalid User Session",
			Error:   401,
		})
	}

	var user models.User
	if err := config.DB.Select("id, name, email").First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(models.Ret{
			Success: false,
			Message: "User not found",
			Error:   404,
		})
	}

	return c.JSON(models.Ret{
		Success: true,
		Message: "User found",
		Error:   200,
		Data:    user,
	})
}

// API Untuk Logout
func Logout(c *fiber.Ctx) error {
	return c.JSON(models.Ret{
		Success: true,
		Message: "Logout successfully",
		Error:   200,
	})
}
