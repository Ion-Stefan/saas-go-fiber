package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Ion-Stefan/saas-go-fiber/config"
	"github.com/Ion-Stefan/saas-go-fiber/database"
	"github.com/Ion-Stefan/saas-go-fiber/internal/handler"
	"github.com/NdoleStudio/lemonsqueezy-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// Connect to the database
	database.ConnectDB()

	// Create a new Fiber instance
	app := fiber.New()

	// Set up middleware
	app.Use(logger.New())

	allowedOrigins := []string{
		config.Envs.WebsiteURL,
		config.Envs.BuildURL,
		config.Envs.LocalURL,
	}

	// CORS configuration
	// Set up CORS
	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Csrf-Token",
		ExposeHeaders:    "X-Csrf-Token",
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			// Check if the incoming origin is in the combined allowed list
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true // Allow the request if the origin matches
				}
			}

			log.Printf("CORS blocked: %s not in allowed origins", origin)
			return false // Block the request if the origin is not allowed
		},
	}))

	// Set up groups for routes
	api := app.Group("/api/v1")
	user := app.Group("/api/user")
	payment := app.Group("/api/payment")

	// Rate limiting middleware
	api.Use(limiter.New(limiter.Config{
		Max:        60,              // Maximum number of requests
		Expiration: 1 * time.Minute, // Time window for the requests
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit based on IP address
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		},
	}))

	// Setup OAuth
	oauthConf := &oauth2.Config{
		ClientID:     config.Envs.ClientID,
		ClientSecret: config.Envs.ClientSecret,
		RedirectURL:  config.Envs.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	paymentClient := lemonsqueezy.New(lemonsqueezy.WithAPIKey(config.Envs.LemonSqueezyWebhookSecret))

	// Link routes
	handler.SetupPaymentRoutes(payment, paymentClient)
	handler.SetupOauthRoutes(user, oauthConf)
	handler.SetupUserRoutes(user)

	// Start the server
	if err := app.Listen(fmt.Sprintf(":%s", config.Envs.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
