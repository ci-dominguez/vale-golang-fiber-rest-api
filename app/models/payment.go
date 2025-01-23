package models

import (
	"github.com/google/uuid"
	"time"
)

// Payment represents a habit entity.
// A payment belongs to a user.
type Payment struct {
	PaymentID   int       `gorm:"primaryKey;autoIncrement" json:"payment_id" db:"payment_id" form:"payment_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id" db:"user_id" form:"user_id"`
	SessionID   string    `gorm:"type:varchar(255);not null" json:"session_id" db:"session_id" form:"session_id"`
	PaymentDate time.Time `gorm:"autoCreateTime" json:"payment_date" db:"payment_date"`
}
