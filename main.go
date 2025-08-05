package main

import (
	"fmt"
	"log"
	"time"

	"baseApi/cache"
	"baseApi/config"
	"baseApi/database"
	"baseApi/logger"
	"baseApi/monitoring"
	"baseApi/routes"
)

/* main is the entry point of the application */
func main() {
	// Initialize logger
	logger.InitLogger()
	logger.Info("Starting CodeBase Golang application...")

	// Load configuration
	cfg := config.LoadConfig()
	logger.Info("Configuration loaded successfully")

	// Initialize database
	database.InitDatabase(cfg)
	logger.Info("Database initialized successfully")

	// Initialize Redis cache
	cache.InitRedis(cfg)
	logger.Info("Redis cache initialized successfully")

	// Initialize Sentry for error tracking
	if cfg.SentryDSN != "" {
		if err := monitoring.InitSentry(cfg); err != nil {
			logger.Error("Failed to initialize Sentry:", err)
		} else {
			logger.Info("Sentry initialized successfully")
			// Ensure Sentry flushes before shutdown
			defer monitoring.FlushSentry(2 * time.Second)
		}
	} else {
		logger.Info("Sentry DSN not provided, skipping Sentry initialization")
	}

	// Setup routes
	router := routes.SetupRoutes()
	logger.Info("Routes setup completed")

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Info("Server starting on port ", cfg.ServerPort)
	
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
