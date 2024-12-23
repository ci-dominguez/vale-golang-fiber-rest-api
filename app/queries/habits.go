package queries

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
	"github.com/google/uuid"
)

// CreateHabit inserts new habit into db
func CreateHabit(habit *models.Habit) error {
	result := database.DB.Create(&habit)

	return result.Error
}

func GetHabitsByUserID(userID uuid.UUID) ([]models.Habit, error) {
	var habits []models.Habit

	result := database.DB.Where("user_id = ?", userID).Find(&habits)
	return habits, result.Error
}

func DeleteHabit(habitID string) error {
	// Convert HabitID string to UUID
	habitUUID, err := uuid.Parse(habitID)
	if err != nil {
		return err
	}

	// Delete habit
	result := database.DB.Where("habit_id = ?", habitUUID).Delete(&models.Habit{})

	return result.Error
}

func VerifyHabitOwnership(habitID string, userUUID uuid.UUID) (bool, error) {
	// Convert habitID string to UUID
	habitUUID, err := uuid.Parse(habitID)
	if err != nil {
		return false, err
	}

	// Check for ownership of habit
	var count int64
	err = database.DB.Model(&models.Habit{}).
		Where("habit_id = ? AND user_id = ?", habitUUID, userUUID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
