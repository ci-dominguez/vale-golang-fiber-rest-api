package queries

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
)

// CreateHabit inserts new habit into db
func CreateHabit(habit *models.Habit) error {
	result := database.DB.Create(&habit)

	return result.Error
}
