package middleware

import (
	"context"
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"time"
)

// Init initializes the Clerk SDK with the secret key.
func Init() {
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
}

// AuthMiddleware is used to authenticate incoming requests.
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("AuthMiddleware triggered for path:", c.Path())

		// Check for auth header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		token := authHeader[len("Bearer "):]
		log.Println("Received Token (middleware):", token)

		// Verify retrieved JWT token
		sessionClaims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
			Token:  token,
			Leeway: 5 * time.Second,
		})
		if err != nil {
			log.Println("JWT Verification Error:", err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		clerkID := sessionClaims.Subject
		log.Println("Extracted clerk ID (middleware):", clerkID)

		var user models.User
		if err := database.DB.Where("clerk_id = ?", clerkID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		c.Locals("userId", user.UserID.String())
		c.Locals("userStatus", user.IsPremium)
		return c.Next()
	}
}
