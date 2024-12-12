package main

import (
	"github.com/ci-dominguez/vale-backend/app/routes"
	"github.com/ci-dominguez/vale-backend/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Env vars
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Init db
	database.InitDB()
	database.MigrateDB()

	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// TO-DO: Register routes
	routes.HabitRoutes(app)

	log.Fatal(app.Listen(":4000"))
}
