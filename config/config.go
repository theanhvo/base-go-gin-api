package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	
	RedisHost     string
	RedisPort     string
	RedisPassword string
	
	ServerPort  string
	JWTSecret   string
	Environment string
	AppVersion  string
	
	// Sentry Configuration
	SentryDSN string
}

/* LoadConfig loads configuration from environment variables */
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "codebase_db"),
		
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-here"),
		Environment: getEnv("ENVIRONMENT", "development"),
		AppVersion:  getEnv("APP_VERSION", "v1.0.0"),
		
		// Sentry
		SentryDSN: getEnv("SENTRY_DSN", ""),
	}
}

/* getEnv gets environment variable with fallback */
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}