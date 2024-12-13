package queries

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
)

// CreateHabitRecords inserts new habit records into db
func CreateHabitRecords(records []models.HabitRecord) error {
	result := database.DB.Create(&records)
	return result.Error
}
