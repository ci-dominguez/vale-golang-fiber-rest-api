package queries

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
	"github.com/google/uuid"
)

// CreateHabit inserts new habit into the db.
func CreateHabit(habit *models.Habit) error {
	result := database.DB.Create(&habit)

	return result.Error
}

// GetHabitsByUserID retrieves all habits associated with the user.
func GetHabitsByUserID(userID uuid.UUID) ([]models.Habit, error) {
	var habits []models.Habit

	result := database.DB.Where("user_id = ?", userID).Find(&habits)
	return habits, result.Error
}

func GetUserHabitCount(userID uuid.UUID) (int64, error) {
	var count int64

	result := database.DB.Model(&models.Habit{}).Where("user_id = ?", userID).Count(&count)

	return count, result.Error
}

// DeleteHabit removes a habit from the db.
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

// UpdateHabit updates a specific habit's name, description, and / or goal.
func UpdateHabit(habitID string, updates map[string]interface{}) (models.Habit, error) {
	var habit models.Habit

	// Convert habitID string to UUID
	habitUUID, err := uuid.Parse(habitID)
	if err != nil {
		return habit, err
	}

	// Validate and sanitize updates map
	validKeys := map[string]bool{"name": true, "description": true, "goal": true}
	sanitizedUpdates := make(map[string]interface{})
	for k, v := range updates {
		if validKeys[k] && v != nil && v != "" {
			sanitizedUpdates[k] = v
		}
	}

	// Update habit
	if err := database.DB.Model(&models.Habit{}).Where("habit_id = ?", habitUUID).Updates(sanitizedUpdates).Error; err != nil {
		return habit, err
	}

	// Fetch and return new habit details5
	if err := database.DB.Where("habit_id = ?", habitUUID).First(&habit).Error; err != nil {
		return habit, err
	}

	return habit, nil
}

// VerifyHabitOwnership checks that the user sending a request is the owner of the habit in the db.
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
