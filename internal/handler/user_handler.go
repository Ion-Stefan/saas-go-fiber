package handler

import (
	"fmt"

	"github.com/Ion-Stefan/saas-go-fiber/config"
	"github.com/Ion-Stefan/saas-go-fiber/internal/service"
	"github.com/Ion-Stefan/saas-go-fiber/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func SetupUserRoutes(app fiber.Router) {
	app.Get("/logout", func(c *fiber.Ctx) error {
		// Clear the cookie
		c.Cookie(&fiber.Cookie{
			Name:     "saas-go-fiber-token",
			Value:    "tokenString",
			Path:     "/",
			SameSite: "Lax",
			HTTPOnly: true,
			Secure:   true,
			MaxAge:   -1,
		})
		return c.SendString("Logout successful")
	})

	// You can add one or more middlewares inline like so
	app.Get("/user_info", middleware.RequireAuth, func(c *fiber.Ctx) error {
		jwtCookie := c.Cookies("saas-go-fiber-token")
		if jwtCookie == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: No token provided")
		}

		token, err := jwt.Parse(jwtCookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.Envs.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))

		userInfo, err := service.GetUserByID(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info: " + err.Error())
		}

		return c.JSON(fiber.Map{
			"user": userInfo,
		})
	})
}
