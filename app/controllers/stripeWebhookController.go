package controllers

import (
	"encoding/json"
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/app/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/webhook"
	"log"
	"os"
	"time"
)

// HandleStripeWebhooks processes incoming webhooks & events from Stripe.
func HandleStripeWebhooks(c *fiber.Ctx) error {
	webhookSecret := os.Getenv("STRIPE_PAYMENT_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("STRIPE_PAYMENT_WEBHOOK_SECRET not configured")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server configuration error",
		})
	}

	// Get request body and signature header
	payload := c.Body()
	sigHeader := c.Get("Stripe-Signature")

	// Construct and verify the incoming event
	event, err := webhook.ConstructEvent(payload, sigHeader, webhookSecret)
	if err != nil {
		log.Println("Webhook verification failed:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Invalid stripe signature")
	}

	switch event.Type {
	case "checkout.session.completed":

		log.Printf("Raw event data: %s\n", string(event.Data.Raw))

		// Extract client_reference_id and id from event.data.object
		var sessionData struct {
			SessionID   string `json:"id"`
			ClientRefID string `json:"client_reference_id"`
		}

		if err := json.Unmarshal(event.Data.Raw, &sessionData); err != nil {
			log.Println("Failed to parse session data:", err)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse stripe payment session data",
			})
		}

		clientReferenceID := sessionData.ClientRefID
		checkoutSessionID := sessionData.SessionID

		// Find user_id from user record where clerk_id = client_reference_id
		userUUID, err := queries.FindUserIDByClerkID(clientReferenceID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve database user from clerk_id",
			})
		}

		// Update is_premium = true for user with user_id found above
		upgradeResult := queries.UpgradeUser(userUUID)
		if upgradeResult != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to upgrade user",
			})
		}

		// Create payment record with client_reference_id, user_id, and timestamp
		var payment models.Payment

		payment.UserID = userUUID
		payment.SessionID = checkoutSessionID
		payment.PaymentDate = time.Now()

		if err := queries.CreatePayment(&payment); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create payment record",
			})
		}

		// Return status code
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User was successfully upgraded",
		})
	default:
		log.Println("Unhandled event type:", event.Type)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Stripe webhook processed successfully",
	})
}
