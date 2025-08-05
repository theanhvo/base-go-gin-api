package database

import (
	"fmt"
	"log"

	"baseApi/config"
	"baseApi/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

/* InitDatabase initializes the database connection */
func InitDatabase(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate the schema
	// if err := AutoMigrate(); err != nil {
	// 	log.Fatal("Failed to migrate database:", err)
	// }
}

/* AutoMigrate runs database migrations */
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
	)
}

/* GetDB returns the database instance */
func GetDB() *gorm.DB {
	return DB
}