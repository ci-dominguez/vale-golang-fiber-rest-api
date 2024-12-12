package database

import (
	"github.com/ci-dominguez/vale-backend/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connection successful")
}

func MigrateDB() {
	err := DB.AutoMigrate(&models.Habit{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration successful")
}
