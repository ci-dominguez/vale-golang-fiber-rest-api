package models

import (
	"time"

	"github.com/google/uuid"
)

// Habit represents a habit entity.
type Habit struct {
	HabitID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"habit_id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name             string    `gorm:"type:varchar(255);not null" json:"name"`
	Description      string    `gorm:"type:text" json:"description,omitempty"`
	TotalCompletions int       `gorm:"default:0" json:"total_completions"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}
