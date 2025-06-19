package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/surahj/ai-mentor-backend/database"
	"github.com/surahj/ai-mentor-backend/routes"
)

func main() {
	// Force load .env file first
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

	// Check if password is set (don't print actual password)
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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	routes.SetupRoutes(e, db)

	// Start server
	fmt.Println("Server starting on :8080")
	log.Fatal(e.Start(":8080"))
}
