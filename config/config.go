package config

import (
	"os"
	"pokemon-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&models.User{})
	DB = db
}

// Config struct to hold the application configuration
type Config struct {
	Port      string
	JwtSecret string
}

// Global variable to store the config instance
var AppConfig *Config

// LoadConfig loads the configuration from environment variables
func LoadConfig() {
	AppConfig = &Config{
		Port:      getEnv("PORT", "8080"),
		JwtSecret: getEnv("JWT_SECRET", "secret"),
	}
}

// Helper function to get environment variables with a fallback default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
