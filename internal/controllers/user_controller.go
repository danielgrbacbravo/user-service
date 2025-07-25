package controllers

import (
	"net/http"

	"github.com/danigrb.dev/auth-service/internal/middleware"
	"github.com/danigrb.dev/auth-service/internal/services"
	"github.com/gin-gonic/gin"
)

// UserController handles user-related routes
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new UserController instance
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// UpdateProfileRequest defines the request body for profile updates
type UpdateProfileRequest struct {
	Email       string         `json:"email,omitempty"`
	Username    string         `json:"username,omitempty"`
	AvatarURL   string         `json:"avatar_url,omitempty"`
	Preferences map[string]any `json:"preferences,omitempty"`
}

// GetProfile handles GET /user/profile
func (uc *UserController) GetProfile(ctx *gin.Context) {
	// Extract user ID from JWT claims using the utility function
	userID, ok := middleware.ExtractUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := uc.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateProfile handles PUT /user/profile
func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	// Extract user ID from JWT claims using the utility function
	userID, ok := middleware.ExtractUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the request to a map for the service
	updates := make(map[string]any)
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.Preferences != nil {
		updates["preferences"] = req.Preferences
	}

	updatedUser, err := uc.userService.UpdateUserProfile(userID, updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles DELETE /user/profile
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	// Extract user ID from JWT claims using the utility function
	userID, ok := middleware.ExtractUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := uc.userService.DeleteUser(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
