package middleware

import (
	"strings"

	"github.com/Ion-Stefan/saas-go-fiber/config"
	"github.com/gofiber/fiber/v2"
)

// RefererCheckMiddleware is a middleware that checks if the request is coming from the allowed referer
func RefererCheckMiddleware(c *fiber.Ctx) error {
	referer := c.Get("Referer")
	allowedReferer := config.Envs.WebsiteURL

	// Check if the referer is empty or if it does not start with the allowed referer
	if referer == "" || !strings.HasPrefix(referer, allowedReferer) {
		return c.Status(fiber.StatusForbidden).SendString("Forbidden from referer")
	}
	// If the referer is allowed, continue
	return c.Next()
}
