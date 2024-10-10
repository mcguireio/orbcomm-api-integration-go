package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	DatabaseURL    string
	S3Bucket       string
	OrbcommAPIKey  string
	OrbcommBaseURL string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		S3Bucket:       os.Getenv("S3_BUCKET"),
		OrbcommAPIKey:  os.Getenv("ORBCOMM_API_KEY"),
		OrbcommBaseURL: os.Getenv("ORBCOMM_BASE_URL"),
	}
}