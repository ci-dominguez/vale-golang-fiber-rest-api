package models

import (
	"time"

	"github.com/google/uuid"
)

// HabitRecord represents a habit record entity.
// A HabitRecord refers to the daily tracking of a Habit.
// E.g., for the month of December 2024, there are 31 HabitRecords.
type HabitRecord struct {
	RecordID    int       `json:"record_id" gorm:"primaryKey;autoIncrement"`
	HabitID     uuid.UUID `json:"habit_id" db:"habit_id"`
	Date        time.Time `json:"date" db:"date"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
}
