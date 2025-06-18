package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/database"
	"github.com/surahj/ai-mentor-backend/routes"
)

func main() {

	// Initialize database
	db, err := database.InitPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Echo
	e := echo.New()

	// Setup routes
	routes.SetupRoutes(e, db)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
