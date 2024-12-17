package queries

import (
	"errors"
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func UpdateHabitRecord(habitID string, recordDate time.Time) (*models.HabitRecord, error) {
	var foundRecord models.HabitRecord

	// Check if record exists
	err := database.DB.Where("habit_id = ? AND date = ?", habitID, recordDate).First(&foundRecord).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Record does not exist, create a new one
		newRecord := models.HabitRecord{
			HabitID:     uuid.MustParse(habitID),
			Date:        recordDate,
			IsCompleted: true,
		}

		if err := database.DB.Create(&newRecord).Error; err != nil {
			return nil, err
		}
		return &newRecord, nil
	} else if err != nil {
		return nil, err
	}

	// Record exists, update it
	foundRecord.IsCompleted = !foundRecord.IsCompleted
	if err := database.DB.Save(&foundRecord).Error; err != nil {
		return nil, err
	}

	return &foundRecord, nil
}
