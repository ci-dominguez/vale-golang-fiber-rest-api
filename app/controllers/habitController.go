package controllers

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/app/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func CreateHabit(c *fiber.Ctx) error {
	var habit models.Habit

	if err := c.BodyParser(&habit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get users id
	userID := c.Locals("userId").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user UUID",
		})
	}

	habit.UserID = userUUID
	habit.HabitID = uuid.New()
	habit.CreatedAt = time.Now()

	if err := queries.CreateHabit(&habit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create habit",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(habit)
}

func GetHabits(c *fiber.Ctx) error {
	userIDInterface := c.Locals("userId")
	println("Retrieved User ID Interface (habitController):", userIDInterface)

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user ID as string",
		})
	}

	println("Retrieved User ID String (habitController):", userIDStr)

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		println("UUID Parse Error:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user UUID",
		})
	}

	// Fetch habits for user
	habits, err := queries.GetHabitsByUserID(userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get habits",
		})
	}

	return c.Status(fiber.StatusOK).JSON(habits)
}
