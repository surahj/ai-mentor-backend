package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/controllers"
	"github.com/surahj/ai-mentor-backend/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Initialize controllers
	userController := controllers.NewUserController(db)


	// Public routes
	e.POST("/api/auth/register", userController.Register)
	e.POST("/api/auth/login", userController.Login)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	// Protected routes
	api := e.Group("/api")
	api.Use(middleware.AuthMiddleware)

	// User routes


	// Learning plan routes


	// Lesson routes
}
