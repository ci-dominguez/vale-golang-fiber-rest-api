package queries

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
)

// CreateUser inserts a new user into the db.
func CreateUser(user *models.User) error {
	result := database.DB.Create(&user)

	return result.Error
}
