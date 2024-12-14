package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUserUUID(c *fiber.Ctx) (uuid.UUID, error) {
	userIDInterface := c.Locals("userId")
	if userIDInterface == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve user ID as string")
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user UUID")
	}

	return userUUID, nil
}
