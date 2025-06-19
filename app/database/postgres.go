package database

import (
	"fmt"
	"os"
	"time"

	"github.com/surahj/ai-mentor-backend/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

func InitPostgres() (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	var db *gorm.DB
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.LearningPlanStructure{},
	)
	if err != nil {
		return nil, err
	}

	dbInstance = db
	return dbInstance, nil
}

func GetDB() *gorm.DB {
	return dbInstance
}
