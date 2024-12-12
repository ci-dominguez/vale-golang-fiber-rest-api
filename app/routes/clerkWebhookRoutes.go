package routes

import (
	"github.com/ci-dominguez/vale-backend/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func ClerkWebhookRoutes(app *fiber.App) {
	api := app.Group("/api/webhooks")
	api.Post("/clerk", controllers.HandleClerkWebhooks)
}
