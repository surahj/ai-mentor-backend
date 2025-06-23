package controllers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/app/models"
	"github.com/surahj/ai-mentor-backend/app/utils"
	"google.golang.org/api/idtoken"
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

	// Check if user already exists
	var existingUser models.User
	result := c.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		// User exists. If they are not verified, we can allow re-sending OTP.
		if existingUser.IsVerified {
			return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "User with this email already exists. Please login to continue.",
			})
		}
	}

	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to process password",
		})
	}

	// Generate OTP
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	otpExpiresAt := time.Now().Add(10 * time.Minute) // OTP valid for 10 minutes

	user := models.User{
		FirstName:       &req.FirstName,
		LastName:        &req.LastName,
		Email:           req.Email,
		Password:        &hashedPass,
		IsVerified:      false,
		OTP:             &otp,
		OTPExpiresAt:    &otpExpiresAt,
		LearningGoal:    req.LearningGoal,
		DailyCommitment: req.DailyCommitment,
		AuthProvider:    "email",
	}

	// If user exists but is not verified, update their record with new OTP.
	// Otherwise, create a new user record.
	if result.Error == nil { // User found
		existingUser.Password = &hashedPass
		existingUser.OTP = &otp
		existingUser.OTPExpiresAt = &otpExpiresAt
		existingUser.FirstName = &req.FirstName
		existingUser.LastName = &req.LastName
		existingUser.LearningGoal = req.LearningGoal
		existingUser.DailyCommitment = req.DailyCommitment
		if err := c.DB.Save(&existingUser).Error; err != nil {
			return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "Failed to update user",
			})
		}
	} else { // User not found, create new
		if err := c.DB.Create(&user).Error; err != nil {
			return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "Failed to create user",
			})
		}
	}

	// Send OTP email
	emailBody := fmt.Sprintf("Hi %s, <br><br>Your One-Time Password (OTP) for AI-Mentor is: <strong>%s</strong>.<br><br>This OTP is valid for 10 minutes. Please use it to complete your registration.<br><br>Thanks,<br>The AI-Mentor Team", *user.FirstName, otp)
	err = c.EmailClient.SendEmail(user.Email, "Your AI-Mentor OTP", emailBody)
	if err != nil {
		log.Printf("Failed to send OTP email to %s: %v", user.Email, err)
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to send verification email.",
		})
	}

	return ctx.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Registration successful. Please check your email for the OTP to verify your account.",
		Data:    nil,
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

	if user.Password == nil || !utils.CheckPasswordHash(loginRequest.Password, *user.Password) {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid credentials",
		})
	}

	if !user.IsVerified {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Account not verified. Please check your email for the OTP.",
		})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to generate token",
		})
	}

	userData := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Login successful",
		Data:    userData,
	})
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

func (c *Controller) VerifyOTP(ctx echo.Context) error {
	var req VerifyOTPRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid request format: " + err.Error(),
		})
	}

	if req.Email == "" || req.OTP == "" {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Email and OTP are required.",
		})
	}

	var user models.User
	if err := c.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "User not found.",
		})
	}

	if user.OTP == nil || *user.OTP != req.OTP {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid OTP.",
		})
	}

	if user.OTPExpiresAt == nil || time.Now().After(*user.OTPExpiresAt) {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "OTP has expired.",
		})
	}

	user.IsVerified = true
	emptyString := ""
	var nilTime *time.Time
	user.OTP = &emptyString
	user.OTPExpiresAt = nilTime
	if err := c.DB.Save(&user).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to verify user.",
		})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to generate token.",
		})
	}

	userData := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Account verified successfully.",
		Data:    userData,
	})
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (c *Controller) ResendOTP(ctx echo.Context) error {
	var req ResendOTPRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid request format: " + err.Error(),
		})
	}

	var user models.User
	if err := c.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "User not found.",
		})
	}

	if user.IsVerified {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "User is already verified.",
		})
	}

	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	otpExpiresAt := time.Now().Add(10 * time.Minute)

	user.OTP = &otp
	user.OTPExpiresAt = &otpExpiresAt
	if err := c.DB.Save(&user).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to update OTP.",
		})
	}

	emailBody := fmt.Sprintf("Your new AI-Mentor OTP is: %s", otp)
	if err := c.EmailClient.SendEmail(user.Email, "Your New AI-Mentor OTP", emailBody); err != nil {
		log.Printf("Failed to resend OTP email to %s: %v", user.Email, err)
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to send OTP email.",
		})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "A new OTP has been sent to your email.",
		Data:    nil,
	})
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	OTP      string `json:"otp" binding:"required,len=6"`
	Password string `json:"password" binding:"required,min=6"`
}

func (c *Controller) ForgotPassword(ctx echo.Context) error {
	var req ForgotPasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid request format: " + err.Error(),
		})
	}

	var user models.User
	if err := c.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "User not found.",
		})
	}

	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	otpExpiresAt := time.Now().Add(10 * time.Minute)

	user.OTP = &otp
	user.OTPExpiresAt = &otpExpiresAt
	if err := c.DB.Save(&user).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to generate reset token.",
		})
	}

	emailBody := fmt.Sprintf("Your password reset OTP is: %s", otp)
	if err := c.EmailClient.SendEmail(user.Email, "Your Password Reset OTP", emailBody); err != nil {
		log.Printf("Failed to send password reset email to %s: %v", user.Email, err)
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to send password reset email.",
		})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Password reset OTP sent to your email.",
		Data:    nil,
	})
}

func (c *Controller) ResetPassword(ctx echo.Context) error {
	var req ResetPasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid request format: " + err.Error(),
		})
	}

	var user models.User
	if err := c.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "User not found.",
		})
	}

	if user.OTP == nil || *user.OTP != req.OTP {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid or expired OTP.",
		})
	}

	if user.OTPExpiresAt == nil || time.Now().After(*user.OTPExpiresAt) {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "OTP has expired.",
		})
	}

	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to process password.",
		})
	}

	user.Password = &hashedPass
	emptyString := ""
	var nilTime *time.Time
	user.OTP = &emptyString
	user.OTPExpiresAt = nilTime
	if err := c.DB.Save(&user).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to reset password.",
		})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Password has been reset successfully.",
		Data:    nil,
	})
}

type GoogleLoginRequest struct {
	Token string `json:"token"`
}

func (c *Controller) GoogleLogin(ctx echo.Context) error {
	var req GoogleLoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid request",
		})
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		log.Println("Google Client ID is not configured")
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "SSO is not configured correctly",
		})
	}

	payload, err := idtoken.Validate(context.Background(), req.Token, googleClientID)
	if err != nil {
		log.Printf("Error validating Google token: %v", err)
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid or expired Google token",
		})
	}

	email := payload.Claims["email"].(string)
	firstName, _ := payload.Claims["given_name"].(string)
	lastName, _ := payload.Claims["family_name"].(string)

	var user models.User
	result := c.DB.Where("email = ?", email).First(&user)

	if result.Error != nil { // User does not exist, create them
		// Generate a random password for Google users
		rand.Seed(time.Now().UnixNano())
		randomPassword := fmt.Sprintf("google_%d_%d", time.Now().Unix(), rand.Intn(1000000))
		hashedPassword, err := utils.HashPassword(randomPassword)
		if err != nil {
			log.Printf("Error hashing password for Google user: %v", err)
			return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "Failed to create user account.",
			})
		}

		newUser := models.User{
			Email:           email,
			FirstName:       &firstName,
			LastName:        &lastName,
			Password:        &hashedPassword,
			IsVerified:      true, // Verified through Google
			AuthProvider:    "google",
			DailyCommitment: 30,
			LearningGoal:    "Not specified",
		}
		if err := c.DB.Create(&newUser).Error; err != nil {
			log.Printf("Error creating user from Google login: %v", err)
			return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "Failed to create user account.",
			})
		}
		user = newUser
	} else { // User exists, just log them in
		// Update auth provider to google if it's not already set
		if user.AuthProvider != "google" {
			user.AuthProvider = "google"
			if err := c.DB.Save(&user).Error; err != nil {
				log.Printf("Error updating auth provider: %v", err)
			}
		}
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		log.Printf("Error generating token for Google user: %v", err)
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to process login.",
		})
	}

	userData := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Login successful",
		Data:    userData,
	})
}
