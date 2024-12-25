package controllers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/svix/svix-webhooks/go"
	"log"
	"net/http"
	"os"
)

// HandleClerkWebhooks processes incoming webhooks & events from Clerk.
func HandleClerkWebhooks(c *fiber.Ctx) error {
	webhookSecret := os.Getenv("CLERK_USER_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("Webhook secret not configured")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server configuration error",
		})
	}

	// Get request body and headers
	payload := c.Body()
	headers := c.GetReqHeaders()

	// Verify the webhook signature
	wh, err := svix.NewWebhook(webhookSecret)
	if err != nil {
		log.Println("Error initializing webhook verification with svix:", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server configuration error",
		})
	}

	var eventPayload map[string]interface{}
	err = wh.Verify(payload, headers)
	if err != nil {
		log.Println("Webhook verification failed:", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid webhook signature",
		})
	}

	// Parse payload into a map
	err = json.Unmarshal(payload, &eventPayload)
	if err != nil {
		log.Println("Error parsing webhook payload:", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload format",
		})
	}

	// Extract event type from payload
	eventType, ok := eventPayload["type"].(string)
	if !ok {
		log.Println("Missing or invalid type in payload")

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event type",
		})
	}

	// Handle events
	switch eventType {
	case "user.created":
		HandleUserCreated(eventPayload)
	default:
		log.Println("Unhandled event type:", eventType)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Webhook processed successfully",
	})
}
