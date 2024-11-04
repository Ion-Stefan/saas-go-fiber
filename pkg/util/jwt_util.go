package util

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetJWTInCookie(c *fiber.Ctx, tokenString string) {
	c.Cookie(&fiber.Cookie{
		Name:     "saas-go-fiber",
		Value:    tokenString,
		Path:     "/",                                 // Available across all paths
		SameSite: "Lax",                               // Allows cookies to be sent with top-level navigations
		HTTPOnly: true,                                // Set to true in production for security
		Expires:  time.Now().Add(time.Hour * 24 * 90), // Set expiration as needed
		Secure:   true,                                // Set to true in production for HTTPS

	})
}
