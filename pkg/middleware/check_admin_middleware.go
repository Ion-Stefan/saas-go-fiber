package middleware

import (
	"github.com/Ion-Stefan/saas-go-fiber/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// Checks for the admin status of the user from the JWT token
func CheckAdminMiddleware(c *fiber.Ctx) error {
	// Get the token from the cookie
	cookie := c.Cookies("saas-go-fiber-token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	// Parse and validate the token
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(config.Envs.JWTSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if isAdmin, ok := claims["admin"].(bool); ok && isAdmin {
			// If the user is an admin, continue
			return c.Next()
		}
	}

	// If the user is not an admin, return an error
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Not an admin",
	})
}
