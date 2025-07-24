package controllers

import (
	"log"
	"net/http"

	"github.com/danigrb.dev/auth-service/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related routes
type AuthController struct {
	userService *services.UserService
}

// NewAuthController creates a new AuthController instance
func NewAuthController() *AuthController {
	return &AuthController{
		userService: services.NewUserService(),
	}
}

// RegisterRequest defines the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest defines the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func (ac *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.userService.CreateUser(req.Email, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token here (will be implemented separately)
	token := "dummy-token" // Placeholder

	ctx.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles user login
func (ac *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.userService.VerifyUserCredentials(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token here
	token := "dummy-token" // Placeholder

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// RefreshToken handles JWT token refresh
func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	// Extract user ID from the token in the Authorization header
	// This would be handled by middleware in a real implementation
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Validate and parse the token
	// This is just a placeholder - actual implementation would verify the token
	log.Println("Received token for refresh:", token)

	// Generate a new token
	newToken := "new-dummy-token" // Placeholder

	ctx.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// AppleLoginRequest defines the request body for Apple login
type AppleLoginRequest struct {
	IdentityToken string `json:"identity_token" binding:"required"`
	UserID        string `json:"user_id" binding:"required"`
	Email         string `json:"email"`
	Username      string `json:"username"`
}

// AppleLogin handles login/registration with Apple credentials
func (ac *AuthController) AppleLogin(ctx *gin.Context) {
	var req AppleLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify Apple token (would be implemented in a real application)
	// For this example, we'll assume the token is valid

	// Use a default username if not provided
	username := req.Username
	if username == "" {
		username = "user_" + req.UserID[:8]
	}

	// Create or get existing user with Apple credentials
	user, err := ac.userService.CreateAppleUser(req.UserID, req.Email, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token := "apple-dummy-token" // Placeholder

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
