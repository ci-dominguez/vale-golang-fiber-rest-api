package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity.
// A user refers to an authenticated and verified user on the app.
// Users are constructed using a ClerkID as a key when added to the db.
// When authenticating the user's ClerkID is used to find and verify user in the db.
type User struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ClerkID   string    `json:"clerk_id" db:"clerk_id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
