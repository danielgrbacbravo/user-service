package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
}

// CreateNewServer initializes the Gin engine, sets up routes, and returns a Server instance.
func CreateNewServer() *Server {
	// Set Gin mode based on environment
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	server := &Server{
		Engine: engine,
	}

	server.SetupRouter()

	return server
}

// Run starts the Gin server on the specified port.
func (s *Server) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting auth-service on port %s", port)
	if err := s.Engine.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
