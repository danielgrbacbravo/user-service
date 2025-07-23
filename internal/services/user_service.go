package services

// User controller placeholder.
// Implement user profile handlers here (e.g., GetProfile, UpdateProfile).
//

import (
	user "github.com/danigrb.dev/auth-service/internal/models"
)

// Temp function to create a new user
func CreateUser(email string, password string) (*user.User, error) {
	u := &user.User{
		Email: email,
	}
	// Assume SetPassword hashes and sets the password
	err := u.SetPassword(password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Temp function to check if a user exists by email
func DoesUserExist(email string) (bool, error) {
	// Placeholder: always returns false, nil
	return false, nil
}

// Temp function to save a user to the database
func SaveUserToDatabase(u *user.User) error {
	// Placeholder: does nothing
	return nil
}

// Temp function to get a user profile by ID
func GetUserProfile(userID uint) (*user.User, error) {
	// Placeholder: returns a dummy user
	return &user.User{
		ID:    userID,
		Email: "dummy@example.com",
	}, nil
}

// Temp function to update a user profile
func UpdateUserProfile(userID uint, newEmail string) (*user.User, error) {
	// Placeholder: returns updated dummy user
	return &user.User{
		ID:    userID,
		Email: newEmail,
	}, nil
}

// Temp function to delete a user by ID
func DeleteUser(userID string) error {
	// Placeholder: does nothing
	return nil
}
