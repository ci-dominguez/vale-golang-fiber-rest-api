package routes

import (
	"github.com/ci-dominguez/vale-backend/app/controllers"
	"github.com/ci-dominguez/vale-backend/app/middleware"
	"github.com/gofiber/fiber/v2"
)

// HabitRoutes registers routes for habit management.
// Routes are protected by AuthMiddleware.
func HabitRoutes(app *fiber.App) {
	api := app.Group("/api/habits")

	api.Use(middleware.AuthMiddleware())

	api.Post("/", controllers.CreateHabit)
	api.Get("/", controllers.GetHabits)
	api.Delete("/", controllers.DeleteHabit)
	api.Patch("/", controllers.UpdateHabit)
}
