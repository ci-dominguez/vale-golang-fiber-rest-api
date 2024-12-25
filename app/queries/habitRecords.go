package queries

import (
	"errors"
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// CreateHabitRecords inserts multiple HabitRecords into the db.
func CreateHabitRecords(records []models.HabitRecord) error {
	result := database.DB.Create(&records)
	return result.Error
}

// GetHabitRecords retrieves HabitRecords for specific Habits within a date range.
func GetHabitRecords(habitIDs []string, startDate time.Time, endDate time.Time) ([]models.HabitRecord, error) {
	var habitRecords []models.HabitRecord

	err := database.DB.Where("habit_id IN ? AND date BETWEEN ? AND ?", habitIDs, startDate, endDate).Find(&habitRecords).Error

	return habitRecords, err
}

// DeleteHabitRecords removes all HabitRecords associated with a specific Habit ID
func DeleteHabitRecords(habitID string) error {
	// Convert habitID string to UUID
	habitUUID, err := uuid.Parse(habitID)
	if err != nil {
		return err
	}

	// Delete habitRecords associated with habit
	result := database.DB.Where("habit_id = ?", habitUUID).Delete(&models.HabitRecord{})

	return result.Error
}

// UpdateHabitRecord toggles or creates a HabitRecord for a specific habit and date.
// If no HabitRecord exists, it creates one with IsCompleted set to true,
// If a HabitRecord exists, it toggles it's IsCompleted field
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
