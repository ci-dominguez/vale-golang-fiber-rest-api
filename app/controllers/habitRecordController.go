package controllers

import (
	"github.com/ci-dominguez/vale-backend/app/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strings"
	"time"
)

func GetHabitRecords(c *fiber.Ctx) error {
	userIDInterface := c.Locals("userId")
	println("Retrieved User ID Interface (habitRecordController):", userIDInterface)

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user ID as string",
		})
	}

	println("Retrieved User ID String (habitRecordController):", userIDStr)

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		println("UUID Parse Error:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user UUID",
		})
	}

	// Validate query params
	habitIDsParam := c.Query("habits")
	dateRangeParam := c.Query("dates")

	if habitIDsParam == "" || dateRangeParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required query param",
		})
	}

	habitIDs := strings.Split(habitIDsParam, ",")

	dateRange := strings.Split(dateRangeParam, ",")
	if len(dateRange) != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date range format",
		})
	}

	startDate, err := time.Parse("2006-01-02", dateRange[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start date format. Expected 'YYYY-MM-DD'",
		})
	}

	endDate, err := time.Parse("2006-01-02", dateRange[1])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid end date format. Expected 'YYYY-MM-DD'",
		})
	}

	if startDate.After(endDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Start date cannot be after end date",
		})
	}

	println("Validated Habit IDs:", habitIDs)
	println("Validated Date Range:", startDate.String(), "to", endDate.String())

	// Verify habits belong to user making the request
	for _, habitID := range habitIDs {
		isAuthorized, err := queries.VerifyHabitOwnership(habitID, userUUID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify habit ownership",
			})
		}

		if !isAuthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You are not authorized to access these habit records",
			})
		}
	}

	println("User is authorized for all requested habit records")

	// Fetch habit records
	habitRecords, err := queries.GetHabitRecords(habitIDs, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch habit records",
			"details": err.Error(),
		})
	}

	println("Number of records fetched: ", len(habitRecords))

	return c.JSON(habitRecords)
}
