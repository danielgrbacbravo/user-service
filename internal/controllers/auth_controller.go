package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danigrb.dev/user-service/internal/middleware"
	"github.com/danigrb.dev/user-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	// Generate JWT token
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":      time.Now().Unix(),                     // Issued at
	}).SignedString(secret)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	// Generate JWT token
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":      time.Now().Unix(),                     // Issued at
	}).SignedString(secret)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// RefreshToken handles JWT token refresh
func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	// Apply the JWTAuth middleware directly to ensure a valid token
	middleware.JWTAuth()(ctx)

	// If the middleware aborted the request, return early
	if ctx.IsAborted() {
		return
	}

	// Extract user ID using the utility function
	userID, ok := middleware.ExtractUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract user ID"})
		return
	}

	// Extract other user information from context set by middleware
	email, _ := ctx.Get("email")
	username, _ := ctx.Get("username")

	// Generate new token with refreshed expiry time
	secret := []byte(os.Getenv("JWT_SECRET"))
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})

	signedToken, err := newToken.SignedString(secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}

// extractBearerToken extracts the JWT token from the Authorization header
func extractBearerToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Missing token")
	}
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}
	return "", fmt.Errorf("Invalid token format")
}

// parseJWTClaims parses and validates the JWT token and returns its claims
func parseJWTClaims(tokenString string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		secret := []byte(os.Getenv("JWT_SECRET"))
		return secret, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("Invalid or expired token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
	}
	return claims, nil
}

// generateJWTToken creates a new JWT token with updated expiry using the provided claims
func generateJWTToken(claims jwt.MapClaims) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  claims["user_id"],
		"email":    claims["email"],
		"username": claims["username"],
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})
	secret := []byte(os.Getenv("JWT_SECRET"))
	return newToken.SignedString(secret)
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
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":      time.Now().Unix(),                     // Issued at
	}).SignedString(secret)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
