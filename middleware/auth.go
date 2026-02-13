package middleware

import (
	"os"
	"strings"

	"github.com/MashuNakamura/todolist-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(401).JSON(models.Ret{
			Success: false,
			Message: "Missing Token",
			Error:   401,
		})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(models.Ret{
			Success: false,
			Message: "Invalid or Expired Token",
			Error:   401,
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(models.Ret{
			Success: false,
			Message: "Invalid Token Claims",
			Error:   401,
		})
	}

	userID := claims["user_id"]
	c.Locals("user_id", userID)
	return c.Next()
}
