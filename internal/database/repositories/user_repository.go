package repositories

import (
	"errors"

	"github.com/danigrb.dev/user-service/internal/database"
	"github.com/danigrb.dev/user-service/internal/database/interfaces"
	"github.com/danigrb.dev/user-service/internal/models"
	"gorm.io/gorm"
)

// Ensure UserRepository implements interfaces.UserRepository
var _ interfaces.UserRepository = (*UserRepository)(nil)

// UserRepository implements the interfaces.UserRepository interface
// using PostgreSQL as the database
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found, but no error
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found, but no error
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found, but no error
		}
		return nil, err
	}
	return &user, nil
}

// FindByAppleID finds a user by Apple ID
func (r *UserRepository) FindByAppleID(appleID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("apple_id = ?", appleID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found, but no error
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user in the database
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user from the database
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// EmailExists checks if an email already exists in the database
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UsernameExists checks if a username already exists in the database
func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
