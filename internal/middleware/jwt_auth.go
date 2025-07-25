package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ExtractUserID extracts the user ID from the Gin context
// Returns the user ID as uint and a boolean indicating success
func ExtractUserID(c *gin.Context) (uint, bool) {
	// Get the user ID from the context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	// Convert the user ID to the expected type
	switch v := userIDValue.(type) {
	case float64:
		return uint(v), true
	case float32:
		return uint(v), true
	case int:
		return uint(v), true
	case uint:
		return v, true
	default:
		return 0, false
	}
}

// JWTAuth is a middleware that validates JWT tokens
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}

		tokenString := authHeader[7:]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secret := []byte(os.Getenv("JWT_SECRET"))
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set claims in context for handlers to use
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("email", claims["email"])
			c.Set("username", claims["username"])
		}

		c.Next()
	}
}
