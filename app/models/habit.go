package models

import (
	"time"

	"github.com/google/uuid"
)

type Habit struct {
    HabitID        uuid.UUID `json:"habit_id" db:"habit_id"`
    UserID         uuid.UUID `json:"user_id" db:"user_id"`
    Name           string    `json:"name" db:"name"`
    Description    string    `json:"description,omitempty" db:"description"`
    TotalCompletions int     `json:"total_completions" db:"total_completions"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
}