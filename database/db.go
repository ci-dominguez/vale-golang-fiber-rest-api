package database

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

// Provides functions to initialize and mange the app's db connection.

// DB is a global variable that holds the db connection instance.
var DB *gorm.DB

// InitDB initializes and establishes a connection to the PostgreSQL db using GORM.
// Logs error if the connection fails.
func InitDB() {
	// Open connection to db on neon
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Assign the db connection instance to DB
	DB = db
	log.Println("Database connection successful")
}

// MigrateDB performs schema migrations for the app's db.
// Logs error if migration fails
func MigrateDB() {
	err := DB.AutoMigrate(&models.User{}, &models.Habit{}, &models.HabitRecord{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration successful")
}
