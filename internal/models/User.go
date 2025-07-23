package models

import "golang.org/x/crypto/bcrypt"

type Preferences map[string]any

type User struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	Email        string      `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash *string     `gorm:"" json:"-"` // Nullable for Apple users
	Username     string      `gorm:"unique;not null" json:"username"`
	AvatarURL    string      `gorm:"" json:"avatar_url,omitempty"`
	Preferences  Preferences `gorm:"type:json" json:"preferences"`
	AppleID      *string     `gorm:"uniqueIndex" json:"-"` // Stores Apple 'sub' claim, nullable for non-Apple users
	AppleEmail   *string     `gorm:"" json:"-"`            // Optional: store Apple email if provided
}

// SetPassword hashes the given password and sets the PasswordHash field.
func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashStr := string(hashed)
	u.PasswordHash = &hashStr
	return nil
}
