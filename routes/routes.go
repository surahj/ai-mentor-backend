package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/controllers"
	"github.com/surahj/ai-mentor-backend/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// User controller
	userController := controllers.NewUserController(db)

	auth := e.Group("/api/auth")
	auth.POST("/register", userController.Register)
	auth.POST("/login", userController.Login)
	auth.GET("/profile", userController.Profile, middleware.AuthMiddleware())
	auth.PUT("/profile", userController.UpdateProfile, middleware.AuthMiddleware())

	// Session controller
	sessionController := controllers.NewSessionController(db)

	sessions := e.Group("/api/sessions")
	sessions.Use(middleware.AuthMiddleware())
	sessions.POST("", sessionController.Create)
	sessions.GET("", sessionController.List)
	sessions.GET("/:id", sessionController.Get)
	sessions.PUT("/:id", sessionController.Update)
	sessions.DELETE("/:id", sessionController.Delete)
	sessions.GET("/tags", sessionController.GetTags)

	// Goal controller
	goalController := controllers.NewGoalController(db)

	goals := e.Group("/api/goals")
	goals.Use(middleware.AuthMiddleware())
	goals.POST("", goalController.Create)
	goals.GET("", goalController.List)
	goals.GET("/:id", goalController.Get)
	goals.PUT("/:id", goalController.Update)
	goals.DELETE("/:id", goalController.Delete)
	goals.GET("/progress", goalController.GetProgress)
}
