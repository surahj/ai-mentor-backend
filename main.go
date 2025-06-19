package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/surahj/ai-mentor-backend/database"
	_ "github.com/surahj/ai-mentor-backend/docs" // This will be auto-generated
	"github.com/surahj/ai-mentor-backend/routes"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title AI Mentor API
// @version 1.0
// @description API for AI Mentor application
// @host localhost:8080
// @BasePath /api

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Debug environment variables
	fmt.Println("=== DEBUG: Environment Variables ===")
	fmt.Printf("DB_HOST: '%s'\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT: '%s'\n", os.Getenv("DB_PORT"))
	fmt.Printf("DB_USER: '%s'\n", os.Getenv("DB_USER"))
	fmt.Printf("DB_NAME: '%s'\n", os.Getenv("DB_NAME"))

	if os.Getenv("DB_PASSWORD") != "" {
		fmt.Println("DB_PASSWORD: '[SET]'")
	} else {
		fmt.Println("DB_PASSWORD: '[MISSING]'")
	}
	fmt.Println("=====================================")

	// Initialize database
	db, err := database.InitPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Database connected successfully!")

	// Initialize Echo
	e := echo.New()

	// Basic middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Simple test route - no middleware, no auth
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// Setup routes
	routes.SetupRoutes(e, db)

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Debug: Print all routes
	fmt.Println("=== Registered Routes ===")
	for _, route := range e.Routes() {
		fmt.Printf("Route: %s %s\n", route.Method, route.Path)
	}
	fmt.Println("=========================")

	// Start server
	fmt.Println("Server starting on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
