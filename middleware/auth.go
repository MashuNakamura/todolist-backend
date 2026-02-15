package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MashuNakamura/todolist-backend/config"
	"github.com/MashuNakamura/todolist-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

	c.Locals("user", token)
	c.Locals("user_id", uint(claims["user_id"].(float64)))
	return c.Next()
}

func getGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin(c *fiber.Ctx) error {
	config := getGoogleConfig()
	url := config.AuthCodeURL("random-state")
	return c.JSON(fiber.Map{
		"url": url,
	})
}

func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	googleConfig := getGoogleConfig()

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to exchange token from Google"})
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to get user info from Google"})
	}
	defer resp.Body.Close()

	userData, _ := io.ReadAll(resp.Body)

	var googleUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		ID    string `json:"id"`
	}
	json.Unmarshal(userData, &googleUser)

	var user models.User
	if err := config.DB.Where("email = ?", googleUser.Email).First(&user).Error; err != nil {
		user = models.User{
			Name:     googleUser.Name,
			Email:    googleUser.Email,
			Password: "",
		}
		if err := config.DB.Create(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Failed to create user"})
		}
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtToken, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to generate token"})
	}

	frontendURL := os.Getenv("FRONTEND_URL")

	return c.Redirect(fmt.Sprintf("%s/auth/google/callback?token=%s", frontendURL, jwtToken))
}
