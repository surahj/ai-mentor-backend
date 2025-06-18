package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/models"
	"github.com/surahj/ai-mentor-backend/utils"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

func (c *UserController) Register(ctx echo.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process password",
		})
	}
	user.Password = hashedPassword

	// Create user
	if err := c.db.Create(&user).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	return ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) Login(ctx echo.Context) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.Bind(&loginRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	var user models.User
	if err := c.db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}
	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
