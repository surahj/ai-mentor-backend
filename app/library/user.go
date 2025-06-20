package library

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/app/database"
	"github.com/surahj/ai-mentor-backend/app/models"
)

func GetUserByID(userID int64) (models.User, error) {

	var user models.User
	db := database.GetDB()

	if err := db.First(&user, userID).Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetUserIDFronContext(ctx echo.Context) (int64, error) {
	userID := ctx.Get("user_id")
	if userID == nil {
		return 0, errors.New("user_id not found")
	}

	switch v := userID.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errors.New("user_id string is not a valid int")
		}
		return id, nil
	default:
		return 0, errors.New("user_id is of unknown type")
	}
}
