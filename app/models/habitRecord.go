package models

import (
	"time"

	"github.com/google/uuid"
)

type HabitRecord struct {
	RecordID    int       `json:"record_id" gorm:"primaryKey;autoIncrement"`
	HabitID     uuid.UUID `json:"habit_id" db:"habit_id"`
	Date        time.Time `json:"date" db:"date"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
}
