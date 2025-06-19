package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"github.com/surahj/ai-mentor-backend/controllers"
	"github.com/surahj/ai-mentor-backend/middleware"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Initialize controllers
	userController := controllers.NewUserController(db)

	auth := e.Group("/api/auth")
	auth.POST("/register", userController.Register)
	auth.POST("/login", userController.Login)

	// Protected route
	auth.GET("/profile", userController.Profile, middleware.AuthMiddleware())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	// User routes


	// Learning plan routes


	// Lesson routes
}
