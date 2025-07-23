package controllers

import (
	"net/http"
	"strconv"

	services "github.com/danigrb.dev/auth-service/internal/services"
	"github.com/gin-gonic/gin"
)

// GetProfile handles GET /user/profile
func GetProfile(ctx *gin.Context) {
	// For demonstration, get userID from query param (in real app, get from JWT or context)
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Placeholder: fetch user profile (replace with DB call)
	u, err := services.GetUserProfile(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user profile"})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

// UpdateProfile handles PUT /user/profile
func UpdateProfile(ctx *gin.Context) {
	// For demonstration, get userID from query param (in real app, get from JWT or context)
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	var req struct {
		Email string `json:"email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Placeholder: update user profile (replace with DB call)
	u, err := services.UpdateUserProfile(uint(userID), req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user profile"})
		return
	}

	ctx.JSON(http.StatusOK, u)
}
