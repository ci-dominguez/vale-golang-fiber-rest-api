package queries

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
)

// CreateUser creates a new user
func CreateUser(user *models.User) error {
	result := database.DB.Create(&user)

	return result.Error
}
