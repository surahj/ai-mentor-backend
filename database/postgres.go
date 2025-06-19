package database

import (
	"fmt"
	"os"

	"github.com/surahj/ai-mentor-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Update your connection string building to include password and SSL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost,
		dbUser,
		dbPassword, // Make sure this is included
		dbName,
		dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Goal{}, // Add this line
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
