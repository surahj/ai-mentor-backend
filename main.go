package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/surahj/ai-mentor-backend/app/configs"
	"github.com/surahj/ai-mentor-backend/app/database"
	app "github.com/surahj/ai-mentor-backend/app/router"
	"github.com/surahj/ai-mentor-backend/docs"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	// Load configuration
	config, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "ai-mentor Service API"
	docs.SwaggerInfo.Description = "This API documents exposes all the available API endpoints for AI Mentor service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("BASE_URL")
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https"}

	ctx := context.Background()
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "dev"
	}

	router := &app.App{}

	// Initialize database
	db, err := database.InitPostgres()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router.Initialize(ctx, db, config)

	router.Run()
}
