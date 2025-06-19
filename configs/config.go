package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}
}

type Config struct {
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}
}

func Load() (*Config, error) {
	var config Config

	// Load from environment variables
	config.DB.Host = os.Getenv("DB_HOST")
	config.DB.User = os.Getenv("DB_USER")
	config.DB.Password = os.Getenv("DB_PASSWORD")
	config.DB.Name = os.Getenv("DB_NAME")
	config.DB.SSLMode = "require" // Default for Supabase

	// Convert port from string to int
	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		portStr = "5432" // Default PostgreSQL port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	config.DB.Port = port

	return &config, nil
}
