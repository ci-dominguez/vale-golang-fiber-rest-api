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
	UserID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"user_id" db:"user_id"`
	ClerkID   string    `gorm:"type:VARCHAR(255);not null" json:"clerk_id" db:"clerk_id"`
	Name      string    `gorm:"type:VARCHAR(255);not null" json:"name" db:"name"`
	Email     string    `gorm:"type:VARCHAR(255);unique;not null" json:"email" db:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" db:"created_at"`
}
