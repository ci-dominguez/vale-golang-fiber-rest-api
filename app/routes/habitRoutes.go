package routes

import (
	"github.com/ci-dominguez/vale-backend/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func HabitRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/habits", controllers.CreateHabit)
}
