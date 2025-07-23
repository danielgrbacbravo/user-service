package server

import (
	"log"
	"os"

	"github.com/danigrb.dev/auth-service/internal/server/routes"

	"github.com/gin-gonic/gin"
)

// Server wraps the Gin engine and configuration.
type Server struct {
	Engine *gin.Engine
}

// NewServer initializes the Gin engine, sets up routes, and returns a Server instance.
func NewServer() *Server {
	// Set Gin mode based on environment
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	// Setup routes
	routes.SetupRoutes(engine)

	return &Server{
		Engine: engine,
	}
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
