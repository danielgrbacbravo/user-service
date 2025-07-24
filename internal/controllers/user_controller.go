package controllers

import (
	"net/http"
	"strconv"

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
	Email       string                 `json:"email,omitempty"`
	Username    string                 `json:"username,omitempty"`
	AvatarURL   string                 `json:"avatar_url,omitempty"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
}

// GetProfile handles GET /user/profile
func (uc *UserController) GetProfile(ctx *gin.Context) {
	// In a real application, you would extract the user ID from the JWT token or session
	// For this example, we'll continue to use the query parameter approach
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

	user, err := uc.userService.GetUserByID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateProfile handles PUT /user/profile
func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	// In a real application, you would extract the user ID from the JWT token or session
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

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the request to a map for the service
	updates := make(map[string]interface{})
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

	updatedUser, err := uc.userService.UpdateUserProfile(uint(userID), updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles DELETE /user/profile
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	// In a real application, you would extract the user ID from the JWT token or session
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

	if err := uc.userService.DeleteUser(uint(userID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
