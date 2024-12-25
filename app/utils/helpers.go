package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Utility functions used across the application

// GetUserUUID retrieves the user's UUID from the Fiber context.
// Expects 'userId' to be stored in the context
// Returns:
// - A valid uuid.UUID is successful
// - An error if:
//  1. The user ID is not found.
//  2. The user ID cannot be cast to a string.
//  3. The string cannot be parsed as a UUID.
func GetUserUUID(c *fiber.Ctx) (uuid.UUID, error) {
	// Get userId from context
	userIDInterface := c.Locals("userId")
	if userIDInterface == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}

	// Cast to string
	userIDStr, ok := userIDInterface.(string)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve user ID as string")
	}

	// Parse into UUID
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user UUID")
	}

	return userUUID, nil
}
