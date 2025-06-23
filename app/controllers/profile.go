package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/app/library"
	"github.com/surahj/ai-mentor-backend/app/models"
)

func (c *Controller) GetProfile(ctx echo.Context) error {
	userID, err := library.GetUserIDFronContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "User not authenticated",
		})
	}

	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "User not found",
		})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Profile retrieved successfully",
		Data: map[string]interface{}{
			"id":                 user.ID,
			"email":              user.Email,
			"first_name":         user.FirstName,
			"last_name":          user.LastName,
			"age":                user.Age,
			"level":              user.Level,
			"background":         user.Background,
			"preferred_language": user.PreferredLanguage,
			"interests":          user.Interests,
			"country":            user.Country,
		},
	})
}

func (c *Controller) UpdateProfile(ctx echo.Context) error {
	var req models.UpdateProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid request format: " + err.Error(),
		})
	}

	userID, err := library.GetUserIDFronContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "User not authenticated",
		})
	}

	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "User not found",
		})
	}

	// Update only the fields that are provided
	if req.Age != nil {
		user.Age = req.Age
	}
	if req.Level != nil {
		user.Level = req.Level
	}
	if req.Background != nil {
		user.Background = req.Background
	}
	if req.PreferredLanguage != nil {
		user.PreferredLanguage = req.PreferredLanguage
	}
	if req.Interests != nil {
		user.Interests = req.Interests
	}
	if req.Country != nil {
		user.Country = req.Country
	}

	if err := c.DB.Save(&user).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to update profile",
		})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Profile updated successfully",
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"id":                 user.ID,
				"email":              user.Email,
				"first_name":         user.FirstName,
				"last_name":          user.LastName,
				"age":                user.Age,
				"level":              user.Level,
				"background":         user.Background,
				"preferred_language": user.PreferredLanguage,
				"interests":          user.Interests,
				"country":            user.Country,
			},
		},
	})
}
