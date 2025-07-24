package database

import (
	"fmt"
	"log"
	"os"

	"github.com/danigrb.dev/auth-service/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectDatabase() {
	LoadEnv()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=public",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Create the schema if it doesn't exist
	db.Exec("CREATE SCHEMA IF NOT EXISTS public")

	DB = db
	fmt.Println("✅ Connected to the database!")

	// Auto migrate models
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to automigrate: ", err)
	}
}
