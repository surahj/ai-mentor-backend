package controllers

import (
	"github.com/surahj/ai-mentor-backend/app/models"
	"github.com/labstack/echo/v4"
)

// RespondJSON makes the error response with payload as json format
func RespondJSON(c echo.Context, code int, message interface{}) error {

	return c.JSON(code, models.ResponseMessage{
		Status:  code,
		Message: message,
	})
}

// RespondRaw makes the error response with payload as json format
func RespondRaw(c echo.Context, code int, message interface{}) error {
	return c.JSON(code, message)
}
