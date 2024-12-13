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

	res := database.DB.Where("user_id = ?", userID).Find(&habits)
	return habits, res.Error
}
