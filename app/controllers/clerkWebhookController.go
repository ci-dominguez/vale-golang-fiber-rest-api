package controllers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/svix/svix-webhooks/go"
	"log"
	"net/http"
	"os"
)

func HandleClerkWebhooks(c *fiber.Ctx) error {
	webhookSecret := os.Getenv("CLERK_USER_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("Webhook secret not configured")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server configuration error",
		})
	}

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

	// Parse payload
	err = json.Unmarshal(payload, &eventPayload)
	if err != nil {
		log.Println("Error parsing webhook payload:", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload format",
		})
	}

	eventType, ok := eventPayload["type"].(string)
	if !ok {
		log.Println("Missing or invalid type in payload")

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event type",
		})
	}

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
