package middleware

import (
	"context"
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

func Init() {
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
}

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	token := authHeader[len("Bearer "):]
	println("Received Token (middleware):", token)

	sessionClaims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
		Token:  token,
		Leeway: 5 * time.Second,
	})
	if err != nil {
		println("JWT Verification Error:", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	clerkID := sessionClaims.Subject
	println("Extracted clerk ID (middleware):", clerkID)

	var user models.User
	if err := database.DB.Where("clerk_id = ?", clerkID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// Store user db id as string
	c.Locals("userId", user.UserID.String())
	return c.Next()
}
