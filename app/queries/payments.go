package queries

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"github.com/ci-dominguez/vale-backend/database"
)

// CreatePayment inserts a new payment record into the payments table.
func CreatePayment(paymentRecord *models.Payment) error {
	result := database.DB.Create(&paymentRecord)

	return result.Error
}
