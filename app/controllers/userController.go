package controllers

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/app/queries"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

// HandleUserCreated processes the "user.created" webhook event from Clerk.
// It:
// 1. Extracts user information (Clerk ID and Email) from the webhook payload
// 2. Generates a new User struct with default values.
// 3. Saves the new user to the db.
func HandleUserCreated(data map[string]interface{}) {
	userData, ok := data["data"].(map[string]interface{})
	if !ok {
		log.Println("Invalid user data format")
		return
	}

	// Getting clerk user idn and email
	clerkID, ok := userData["id"].(string)
	if !ok {
		log.Println("Clerk IDN not found or invalid")
		return
	}

	emailAddresses, ok := userData["email_addresses"].([]interface{})
	if !ok || len(emailAddresses) == 0 {
		log.Println("Email address not found or empty")
		return
	}

	email, ok := emailAddresses[0].(map[string]interface{})["email_address"].(string)
	if !ok {
		log.Println("Email address not found")
		return
	}

	name := strings.Split(email, "@")[0]

	// Create user
	user := models.User{
		UserID:    uuid.New(),
		ClerkID:   clerkID,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	// Save user in db
	if err := queries.CreateUser(&user); err != nil {
		log.Printf("Failed to create user: %v\n", err)
	} else {
		log.Printf("Successfully created user: %v\n", user)
	}
}
