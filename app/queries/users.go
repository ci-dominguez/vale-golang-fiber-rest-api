package queries

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
	"github.com/google/uuid"
)

// CreateUser inserts a new user into the db.
func CreateUser(user *models.User) error {
	result := database.DB.Create(&user)

	return result.Error
}

// FindUserIDByClerkID retrieves the database user_id based on the authenticated clerk_id.
func FindUserIDByClerkID(clerkID string) (uuid.UUID, error) {
	var user models.User

	result := database.DB.Where("clerk_id = ?", clerkID).First(&user)

	return user.UserID, result.Error
}

// UpgradeUser updates the is_premium column of the user to true.
func UpgradeUser(userID uuid.UUID) error {

	result := database.DB.Model(&models.User{}).Where("user_id = ?", userID).Update("is_premium", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
