package controllers

import (
	"fmt"
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/app/queries"
	"github.com/ci-dominguez/vale-backend/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// CreateHabit handles the creation of a new Habit.
// It:
// 1. Parses request body into a Habit struct
// 2. Retrieves the authenticated user's UUID from the request context.
// 3. Creates a new Habit in the db with a month's worth of HabitRecords.
func CreateHabit(c *fiber.Ctx) error {
	var habit models.Habit

	// Get users db id
	userUUID, err := utils.GetUserUUID(c)
	if err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get users premium status.
	isUserPremium, err := utils.GetUserStatus(c)
	if err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// If user is not premium, check habit limits.
	if !isUserPremium {

		habitCount, err := queries.GetUserHabitCount(userUUID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get habit count",
			})
		}

		if habitCount >= 3 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User has reached free tier limits.",
			})
		}
	}

	// Parse JSON into habit struct
	if err := c.BodyParser(&habit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	println("Parsed Habit Name:", habit.Name)
	println("Parsed Habit Description:", habit.Description)

	// Assign non-user generated values to habit
	habit.UserID = userUUID
	habit.HabitID = uuid.New()
	habit.CreatedAt = time.Now()

	// Randomly assign a color to the habit
	colors := [...]string{
		"purple",
		"green",
		"red",
		"orange",
		"yellow",
		"light-yellow",
		"lime",
		"light-lime",
		"mint",
		"light-mint",
		"teal",
		"light-cyan",
		"indigo",
		"light-indigo",
		"violet",
		"pink",
		"light-pink",
		"rose",
	}

	randIndex := rand.Intn(len(colors))

	randColor := colors[randIndex]

	habit.Color = randColor

	if err := queries.CreateHabit(&habit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create habit",
		})
	}

	// Create a month worth of habitRecords for the habit
	startDate := time.Now()
	endDate := startDate.AddDate(0, 1, 0)

	var habitRecords []models.HabitRecord
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		habitRecords = append(habitRecords, models.HabitRecord{
			HabitID:     habit.HabitID,
			Date:        d,
			IsCompleted: false,
		})
	}

	if err := queries.CreateHabitRecords(habitRecords); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create default habit records",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(habit)
}

// GetHabits retrieves all habits associated with the authenticated user.
// It:
// 1. Retrieves the user's UUID from the request context.
// 2. Fetches all habits belonging to that user from the db.
func GetHabits(c *fiber.Ctx) error {
	// Get users db id
	userUUID, err := utils.GetUserUUID(c)
	if err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
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

// DeleteHabit deletes a specific Habit and its associated HabitRecords.
// It:
// 1. Retrieves the user's UUID from the request context.
// 2. Validates that the required query param is provided.
// 3. Verifies that the authenticated user owns the specified Habit.
// 4. Deletes all associated HabitRecords before deleting the Habit itself.
func DeleteHabit(c *fiber.Ctx) error {
	// Get users db id
	userUUID, err := utils.GetUserUUID(c)
	if err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate query params
	habitID := c.Query("habit")

	if habitID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required query param",
		})
	}

	println("Parsed Habit ID:", habitID)

	// Verify habit belongs to user making the request
	isAuthorized, err := queries.VerifyHabitOwnership(habitID, userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify habit ownership",
		})
	}

	if !isAuthorized {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	println("User is authorized to modify habit")

	// Delete all habitRecords for habit
	if err := queries.DeleteHabitRecords(habitID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete habit records",
		})
	}

	// Delete habit
	if err := queries.DeleteHabit(habitID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete habit",
		})
	}

	println("Successfully deleted habit and all its records.")

	return c.Status(fiber.StatusOK).JSON(habitID)
}

// UpdateHabit updates specific details of a Habit.
func UpdateHabit(c *fiber.Ctx) error {
	// Get users db id
	userUUID, err := utils.GetUserUUID(c)
	if err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Parse request body into a map
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	fmt.Printf("Parsed updates: %+v\n", updates)

	// Validate habitID from URL params
	habitID := c.Query("habit")
	if habitID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required query param: habit",
		})
	}

	println("Parsed Habit ID:", habitID)

	// Verify habit belongs to user making the request
	isAuthorized, err := queries.VerifyHabitOwnership(habitID, userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify habit ownership",
		})
	}
	if !isAuthorized {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	println("User is authorized to modify habit")

	// Update habit
	habit, err := queries.UpdateHabit(habitID, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update habit",
		})
	}

	return c.Status(fiber.StatusOK).JSON(habit)
}
