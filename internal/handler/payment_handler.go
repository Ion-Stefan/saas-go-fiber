package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/Ion-Stefan/saas-go-fiber/config"
	"github.com/NdoleStudio/lemonsqueezy-go"
	"github.com/gofiber/fiber/v2"
)

// SetupLemonSqueezyRoutes configures the Fiber router to handle LemonSqueezy webhook events
func SetupPaymentRoutes(app fiber.Router, client *lemonsqueezy.Client) {
	// Create a LemonSqueezy client instance

	app.Post("/webhook", func(c *fiber.Ctx) error {
		// Retrieve the raw body and the LemonSqueezy X-Signature header
		payload := c.Body()
		sigHeader := c.Get("X-Signature")

		// Verify the webhook signature
		if !verifySignature(payload, sigHeader, config.Envs.LemonSqueezyWebhookSecret) {
			log.Printf("Webhook signature verification failed")
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid webhook signature.",
			})
		}

		// Retrieve the event type from the header
		eventName := c.Get("X-Event-Name")

		// Handle the event based on its type
		switch eventName {
		case lemonsqueezy.WebhookEventOrderCreated:
			// Handle order created event
			var order lemonsqueezy.WebhookRequestOrder
			if err := json.Unmarshal(payload, &order); err != nil {
				log.Printf("Error parsing order: %v\n", err)
				return c.Status(400).JSON(fiber.Map{
					"error": "Failed to parse order data.",
				})
			}
			// Implement your order logic here
			balance := int64(order.Data.Attributes.Subtotal)
			log.Printf("Order created with ID: %s for amount: %d\n", order.Data.Attributes.UserEmail, balance)

		case lemonsqueezy.WebhookEventOrderRefunded:
			// Handle order refunded event
			var refund lemonsqueezy.WebhookRequestOrder
			if err := json.Unmarshal(payload, &refund); err != nil {
				log.Printf("Error parsing refund: %v\n", err)
				return c.Status(400).JSON(fiber.Map{
					"error": "Failed to parse refund data.",
				})
			}
			// Implement your refund logic here
			balance := int64(refund.Data.Attributes.Subtotal)
			log.Printf("Refund created with ID: %s for amount: %d\n", refund.Data.Attributes.UserEmail, balance)

		default:
			// Handle unrecognized event types
			log.Printf("Unhandled event type: %s\n", eventName)
		}

		// Return a success response to LemonSqueezy
		return c.SendStatus(200)
	})
}

// verifySignature verifies the LemonSqueezy webhook signature using HMAC-SHA256
func verifySignature(payload []byte, signatureHeader, secret string) bool {
	// Create a new HMAC using SHA256 and the webhook secret
	h := hmac.New(sha256.New, []byte(secret))

	// Write the payload to the HMAC object
	h.Write(payload)

	// Compute the HMAC digest and encode it as a hexadecimal string
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Compare the generated signature with the signature from the header
	return hmac.Equal([]byte(expectedSignature), []byte(signatureHeader))
}
