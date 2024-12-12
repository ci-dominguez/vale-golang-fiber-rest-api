package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
    UserID    uuid.UUID `json:"user_id" db:"user_id"`
    ClerkID   string    `json:"clerk_id" db:"clerk_id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}