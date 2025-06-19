package library

import (
	"log"

	"github.com/surahj/ai-mentor-backend/app/database"
	"github.com/surahj/ai-mentor-backend/app/models"
)

func GetUserByID(userID int64) (models.User, error) {

	var user models.User
	db := database.GetDB()

	log.Printf("Getting user by ID: %d", userID)

	if err := db.First(&user, userID).Error; err != nil {
		return user, err
	}
	
	return user, nil
}
