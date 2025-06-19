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

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uc *UserController) Register(c echo.Context) error {
	var input RegisterInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Make sure email is not empty
	if input.Email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email is required"})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
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

func (uc *UserController) Profile(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	var user models.User
	if err := uc.db.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":    user.ID,
		"email": user.Email,
		// Add other fields as needed
	})
}
