package main

import (
	"log"
	"time"

	"orbcomm-ship-tracker/api"
	"orbcomm-ship-tracker/config"
	"orbcomm-ship-tracker/database"
	"orbcomm-ship-tracker/orbcomm"
	"orbcomm-ship-tracker/s3"
	"orbcomm-ship-tracker/tasks"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Initialize S3 client
	s3Client, err := s3.NewS3Client()
	if err != nil {
		logger.Fatal("Failed to initialize S3 client", zap.Error(err))
	}

	// Initialize Orbcomm client
	orbcommClient := orbcomm.NewClient(cfg.OrbcommBaseURL, cfg.OrbcommAPIKey)

	// Initialize Echo framework
	e := echo.New()

	// Set up API routes
	api.SetupRoutes(e, db, s3Client, orbcommClient, logger)

	// Start background tasks
	go tasks.StartOrbcommPoller(logger, time.Minute*5, db, orbcommClient, s3Client, cfg.S3Bucket)

	// Start the server
	logger.Info("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}