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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
