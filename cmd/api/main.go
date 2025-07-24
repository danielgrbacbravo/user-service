package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/danigrb.dev/auth-service/internal/server"
)

func main() {
	// Load environment variables from .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Gin router with all routes configured
	server := server.CreateNewServer()

	// Get port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting auth-service on port %s", port)
	server.Run()
}
