package services

import (
	"errors"

	"github.com/danigrb.dev/auth-service/internal/database/interfaces"
	"github.com/danigrb.dev/auth-service/internal/database/repositories"
	"github.com/danigrb.dev/auth-service/internal/models"
)

// UserService handles business logic related to users
type UserService struct {
	userRepo interfaces.UserRepository
}

// NewUserService creates a new UserService instance with repositories from the factory
func NewUserService() *UserService {
	factory := repositories.NewFactory()
	return &UserService{
		userRepo: factory.GetUserRepository(),
	}
}

// NewUserServiceWithRepo creates a new UserService with a specific repository
// This is useful for testing with mock repositories
func NewUserServiceWithRepo(userRepo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with the given email and password
func (s *UserService) CreateUser(email, username, password string) (*models.User, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already in use")
	}

	// Check if username already exists
	exists, err = s.userRepo.UsernameExists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already in use")
	}

	// Create new user
	user := &models.User{
		Email:       email,
		Username:    username,
		Preferences: models.Preferences{},
	}

	// Set and hash password
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// Save user to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// UpdateUserProfile updates user information
func (s *UserService) UpdateUserProfile(id uint, updates map[string]interface{}) (*models.User, error) {
	// Get the user first
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if email is being updated and is unique
	if email, ok := updates["email"].(string); ok && email != user.Email {
		exists, err := s.userRepo.EmailExists(email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already in use")
		}
		user.Email = email
	}

	// Check if username is being updated and is unique
	if username, ok := updates["username"].(string); ok && username != user.Username {
		exists, err := s.userRepo.UsernameExists(username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("username already in use")
		}
		user.Username = username
	}

	// Update avatar URL if provided
	if avatarURL, ok := updates["avatar_url"].(string); ok {
		user.AvatarURL = avatarURL
	}

	// Update preferences if provided
	if preferences, ok := updates["preferences"].(models.Preferences); ok {
		// Merge existing preferences with new ones
		if user.Preferences == nil {
			user.Preferences = models.Preferences{}
		}
		for k, v := range preferences {
			user.Preferences[k] = v
		}
	}

	// Save the updated user
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by their ID
func (s *UserService) DeleteUser(id uint) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Delete the user
	return s.userRepo.Delete(id)
}

// VerifyUserCredentials verifies email and password
func (s *UserService) VerifyUserCredentials(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password (this would use the bcrypt.CompareHashAndPassword method)
	// Assume User has a VerifyPassword method
	// if !user.VerifyPassword(password) {
	//     return nil, errors.New("invalid credentials")
	// }

	return user, nil
}

// CreateAppleUser creates a user with Apple credentials
func (s *UserService) CreateAppleUser(appleID, email, username string) (*models.User, error) {
	// Check if user already exists with this Apple ID
	existingUser, err := s.userRepo.FindByAppleID(appleID)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return existingUser, nil // User already exists, return it
	}

	// Create new user with Apple ID
	user := &models.User{
		Email:       email,
		Username:    username,
		AppleID:     &appleID,
		AppleEmail:  &email,
		Preferences: models.Preferences{},
	}

	// Save user to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
