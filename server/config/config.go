package config

import (
	"os"
	"path/filepath"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Media    MediaConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int
	Host string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path string
}

// MediaConfig holds media-related configuration
type MediaConfig struct {
	Paths map[string]string // Type -> Path mapping
}

// LoadConfig loads configuration from environment variables or defaults
func LoadConfig() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("FINDARR_PORT", 8080),
			Host: getEnvAsString("FINDARR_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Path: getEnvAsString("FINDARR_DB_PATH", filepath.Join("config", "findarr.db")),
		},
		Media: MediaConfig{
			Paths: make(map[string]string),
		},
	}

	// Add default media paths
	cfg.Media.Paths["movies"] = getEnvAsString("FINDARR_MOVIES_PATH", "")
	cfg.Media.Paths["shows"] = getEnvAsString("FINDARR_SHOWS_PATH", "")
	cfg.Media.Paths["books"] = getEnvAsString("FINDARR_BOOKS_PATH", "")
	cfg.Media.Paths["music"] = getEnvAsString("FINDARR_MUSIC_PATH", "")

	return cfg
}

// Helper function to get an environment variable as a string with a default value
func getEnvAsString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get an environment variable as an integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

