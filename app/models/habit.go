package models

import (
	"time"

	"github.com/google/uuid"
)

// Habit represents a habit entity.
// A habit belongs to a user.
type Habit struct {
	HabitID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"habit_id" db:"habit_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id" db:"user_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name" db:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty" db:"description"`
	Goal        int       `gorm:"default:0" json:"goal" db:"goal"`
	Color       string    `gorm:"type:varchar(255)" json:"color" db:"color"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at" db:"created_at"`
}
