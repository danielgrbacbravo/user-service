package server

import (

	//	"github.com/danigrb.dev/auth-service/internal/middleware"

	"github.com/danigrb.dev/auth-service/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin engine and sets up all routes and middleware.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Add global middleware here (e.g., logging, CORS)
	// TODO: Uncomment and implement middleware when available
	// router.Use(middleware.Logging())
	// router.Use(middleware.CORS())

	// Auth routes
	auth := router.Group("/auth")
	{
		// TODO: Update handler signatures in controllers to use func(ctx *gin.Context)
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/apple", controllers.AppleLogin)
	}

	// User profile routes (protected)
	// TODO: Implement user profile routes when controllers and middleware are available
	// user := router.Group("/user")
	// user.Use(middleware.JWTAuth())
	// {
	// 	user.GET("/profile", controllers.GetProfile)
	// 	user.PUT("/profile", controllers.UpdateProfile)
	// }

	// TODO: Implement health check route when controllers.HealthCheck is available
	// router.GET("/health", controllers.HealthCheck)

	return router
}
