package server

import (
	"github.com/danigrb.dev/auth-service/internal/controllers"
	"github.com/danigrb.dev/auth-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin engine and sets up all routes and middleware.
func (server *Server) SetupRouter() {
	router := server.Engine

	// Add global middleware here (e.g., logging, CORS)
	// TODO: Uncomment and implement middleware when available
	// router.Use(middleware.Logging())
	// router.Use(middleware.CORS())

	// Initialize controllers
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.RefreshToken)
		auth.POST("/apple", authController.AppleLogin)
	}

	// User profile routes
	user := router.Group("/user")
	user.Use(middleware.JWTAuth())
	{
		user.GET("/profile", userController.GetProfile)
		user.PUT("/profile", userController.UpdateProfile)
		user.DELETE("/profile", userController.DeleteUser)
	}

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
