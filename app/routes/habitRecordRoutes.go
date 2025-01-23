package routes

import (
	"github.com/ci-dominguez/vale-backend/app/controllers"
	"github.com/ci-dominguez/vale-backend/app/middleware"
	"github.com/gofiber/fiber/v2"
)

// HabitRecordRoutes registers routes for managing habitRecords.
// Routes are protected by AuthMiddleware.
func HabitRecordRoutes(app *fiber.App) {
	api := app.Group("/api/habit-records")

	api.Use(middleware.AuthMiddleware())

	api.Get("/", controllers.GetHabitRecords)
	api.Patch("/", controllers.UpdateHabitRecord)
}
