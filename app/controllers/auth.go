package controllers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/app/models"
	"github.com/surahj/ai-mentor-backend/app/utils"
)

func (c *Controller) Register(ctx echo.Context) error {
	var req models.CreateUserRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: err.Error(),
		})
	}

	if req.Email == "" {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Email is required",
		})
	}

	if req.Password == "" {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Password is required",
		})
	}

	if req.DailyCommitment == 0 {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Daily commitment is required",
		})
	}

	if req.LearningGoal == "" {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Learning goal is required",
		})
	}

	// check if user already exists
	var existingUser models.User
	if err := c.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "User already exists",
		})
	}

	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to process password",
		})
	}

	// Create user
	user := models.User{
		Email:           req.Email,
		Password:        hashedPass,
		FirstName:       &req.FirstName,
		LastName:        &req.LastName,
		DailyCommitment: req.DailyCommitment,
		LearningGoal:    req.LearningGoal,
	}

	if err := c.DB.Create(&user).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to generate token",
		})
	}

	// Create user response data with token
	userData := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       *user.FirstName,
			"last_name":        *user.LastName,
			"daily_commitment": user.DailyCommitment,
			"learning_goal":    user.LearningGoal,
		},
	}

	return ctx.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "User created successfully",
		Data:    userData,
	})
}

func (c *Controller) Login(ctx echo.Context) error {
	var loginRequest models.LoginRequest

	if err := ctx.Bind(&loginRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: err.Error(),
		})
	}

	if loginRequest.Email == "" {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Email is required",
		})
	}

	if loginRequest.Password == "" {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Password is required",
		})
	}

	var user models.User
	if err := c.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid credentials",
		})
	}

	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid credentials",
		})
	}
	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)

	if err != nil {
		log.Printf("Error generating token: %v", err)
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to generate token",
		})
	}

	// Create user response data
	userData := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":               user.ID,
			"email":            user.Email,
			"first_name":       *user.FirstName,
			"last_name":        *user.LastName,
			"daily_commitment": user.DailyCommitment,
			"learning_goal":    user.LearningGoal,
		},
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Login successful",
		Data:    userData,
	})
}
