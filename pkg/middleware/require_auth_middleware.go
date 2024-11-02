package middleware

import (
	"github.com/Ion-Stefan/saas-go-fiber/config"
	"github.com/Ion-Stefan/saas-go-fiber/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// RequireAuth is a middleware that checks if the user is authenticated
func RequireAuth(c *fiber.Ctx) error {
	// Get the token from the cookie
	cookie := c.Cookies("saas-go-fiber-token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse the token and check if it is valid
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Envs.JWTSecret), nil
	})
	// If the token is invalid, return an error
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	// Check the user id
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	_, err = service.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting user",
		})
	}
	// If the token is valid, continue
	return c.Next()
}
