package routes

import (
	"github.com/ci-dominguez/vale-backend/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// StripeWebhookRoutes registers routes to handle incoming webhook events from Stripe.
func StripeWebhookRoutes(app *fiber.App) {
	api := app.Group("/api/webhooks/stripe")
	api.Post("/", controllers.HandleStripeWebhooks)
}
