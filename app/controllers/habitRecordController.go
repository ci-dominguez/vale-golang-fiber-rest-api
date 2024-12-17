package controllers

import (
	"github.com/ci-dominguez/vale-backend/app/queries"
	"github.com/ci-dominguez/vale-backend/app/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

func GetHabitRecords(c *fiber.Ctx) error {
	// Get users db id
	userUUID, err := utils.GetUserUUID(c)
	if err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
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

func UpdateHabitRecord(c *fiber.Ctx) error {

	// Get users db id
	userUUID, err := utils.GetUserUUID(c)
	if err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate query params
	habitID := c.Query("habit")
	recordDateParam := c.Query("date")

	if habitID == "" || recordDateParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required query param",
		})
	}

	recordDate, err := time.Parse("2006-01-02", recordDateParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid record date format. Expected 'YYYY-MM-DD'",
		})
	}

	println("Habit ID:", habitID)
	println("Habit Record Date:", recordDate.String())

	// Verify habit belongs to user making the request
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

	// Update habit record
	updatedHabitRecord, err := queries.UpdateHabitRecord(habitID, recordDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update habit record",
		})
	}

	println("Updated Habit Record:", updatedHabitRecord)

	return c.JSON(updatedHabitRecord)
}
