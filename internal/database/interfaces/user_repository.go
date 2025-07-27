package interfaces

import (
	"github.com/danigrb.dev/user-service/internal/models"
)

// UserRepository defines the interface for user database operations
type UserRepository interface {
	// Create a new user
	Create(user *models.User) error

	// Find a user by ID
	FindByID(id uint) (*models.User, error)

	// Find a user by email
	FindByEmail(email string) (*models.User, error)

	// Find a user by username
	FindByUsername(username string) (*models.User, error)

	// Find a user by Apple ID
	FindByAppleID(appleID string) (*models.User, error)

	// Update a user
	Update(user *models.User) error

	// Delete a user
	Delete(id uint) error

	// Check if email exists
	EmailExists(email string) (bool, error)

	// Check if username exists
	UsernameExists(username string) (bool, error)
}
