package controllers

import (
	"encoding/json"
	//"github.com/ci-dominguez/vale-backend/app/models"
	//"github.com/ci-dominguez/vale-backend/app/queries"
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

func HandleUserCreated(data map[string]interface{}) {
	userData, ok := data["data"].(map[string]interface{})
	if !ok {
		log.Println("Invalid user data format")
		return
	}

	// Getting clerk user idn and email
	clerkID, ok := userData["id"].(string)
	if !ok {
		log.Println("Clerk IDN not found or invalid")
		return
	}

	emailAddresses, ok := userData["email_addresses"].([]interface{})
	if !ok || len(emailAddresses) == 0 {
		log.Println("Email address not found or empty")
		return
	}

	email, ok := emailAddresses[0].(map[string]interface{})["email_address"].(string)
	if !ok {
		log.Println("Email address not found")
		return
	}

	log.Printf("Clerk ID: %s Email: %s\n", clerkID, email)
}
