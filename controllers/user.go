package controllers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/models"
	"github.com/surahj/ai-mentor-backend/utils"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

// Create a validator once when package initializes
var validate = validator.New()

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

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateProfileInput struct {
	Name            string `json:"name" validate:"required"`
	LearningGoal    string `json:"learning_goal" validate:"required"`
	DailyCommitment int    `json:"daily_commitment" validate:"gte=0,lte=1440"` // Max minutes in a day
}

// UpdateProfile handles updating user profile information
func (uc *UserController) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	var input UpdateProfileInput
	if err := validateInput(c, &input); err != nil {
		return err
	}

	// Get user from database
	var user models.User
	if err := uc.db.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	// Update fields
	user.Name = input.Name
	user.LearningGoal = input.LearningGoal
	user.DailyCommitment = input.DailyCommitment

	// Save changes
	if err := uc.db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update profile"})
	}

	// Return updated user info
	return c.JSON(http.StatusOK, echo.Map{
		"id":               user.ID,
		"email":            user.Email,
		"name":             user.Name,
		"learning_goal":    user.LearningGoal,
		"daily_commitment": user.DailyCommitment,
	})
}

// Utility function to validate struct
func validateInput(c echo.Context, input interface{}) error {
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input format"})
	}

	if err := validate.Struct(input); err != nil {
		// Return validation errors
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	return nil
}

// For requesting password reset
type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

// For resetting password
type ResetPasswordInput struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

// Generate a random token
func generateToken() string {
	// We'll use a combination of time and crypto/rand for simplicity
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", b)
}

// Request a password reset
func (uc *UserController) ForgotPassword(c echo.Context) error {
	var input ForgotPasswordInput
	if err := validateInput(c, &input); err != nil {
		return err
	}

	// Check if user exists
	var user models.User
	if err := uc.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Don't reveal if email exists or not for security
		return c.JSON(http.StatusOK, echo.Map{
			"message": "If your email exists in our system, you'll receive a reset link shortly",
		})
	}

	// Generate token and expiry
	token := generateToken()
	expiry := time.Now().Add(1 * time.Hour)

	// Update user with token
	user.ResetToken = token
	user.ResetTokenExpires = expiry
	if err := uc.db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to process request"})
	}

	// Create reset link
	resetURL := fmt.Sprintf("https://yourdomain.com/reset-password?token=%s", token)

	// Send email (this will just log in development mode)
	subject := "Password Reset Request"
	body := fmt.Sprintf("Please use this link to reset your password: %s\nThis link will expire in 1 hour.", resetURL)

	if err := utils.SendEmail(user.Email, subject, body); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to send email"})
	}

	response := echo.Map{
		"message": "If your email exists in our system, you'll receive a reset link shortly",
	}

	// For development, include the token in the response
	// REMOVE THIS IN PRODUCTION!
	if os.Getenv("EMAIL_FROM") == "your-email@gmail.com" {
		response["dev_token"] = token
		response["dev_reset_url"] = resetURL
	}

	return c.JSON(http.StatusOK, response)
}

// Reset password with token
func (uc *UserController) ResetPassword(c echo.Context) error {
	var input ResetPasswordInput
	if err := validateInput(c, &input); err != nil {
		return err
	}

	// Find user with token
	var user models.User
	if err := uc.db.Where("reset_token = ?", input.Token).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or expired token"})
	}

	// Check if token is expired
	if user.ResetTokenExpires.Before(time.Now()) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Token has expired"})
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to process password"})
	}

	// Update user
	user.Password = hashedPassword
	user.ResetToken = ""
	user.ResetTokenExpires = time.Time{}

	if err := uc.db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update password"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Password has been reset successfully"})
}
