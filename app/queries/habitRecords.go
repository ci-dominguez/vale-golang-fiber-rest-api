package queries

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
	"time"
)

// CreateHabitRecords inserts new habit records into db
func CreateHabitRecords(records []models.HabitRecord) error {
	result := database.DB.Create(&records)
	return result.Error
}

func GetHabitRecords(habitIDs []string, startDate time.Time, endDate time.Time) ([]models.HabitRecord, error) {
	var habitRecords []models.HabitRecord

	err := database.DB.Where("habit_id IN ? AND date BETWEEN ? AND ?", habitIDs, startDate, endDate).Find(&habitRecords).Error

	return habitRecords, err
}
