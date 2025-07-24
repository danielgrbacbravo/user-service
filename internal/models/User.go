package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Preferences map[string]any

// Scan implements the sql.Scanner interface for Preferences.
// This allows GORM to properly scan JSON data from the database into Preferences.
func (p *Preferences) Scan(src any) error {
	if src == nil {
		*p = make(Preferences)
		return nil
	}

	// Type assertion to get the bytes
	var source []byte
	switch src := src.(type) {
	case string:
		source = []byte(src)
	case []byte:
		source = src
	default:
		return errors.New("incompatible type for Preferences")
	}

	// Unmarshal JSON data
	var result map[string]any
	if err := json.Unmarshal(source, &result); err != nil {
		return err
	}

	*p = result
	return nil
}

// Value implements the driver.Valuer interface for Preferences.
// This allows GORM to properly store Preferences as JSON in the database.
func (p Preferences) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

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

// VerifyPassword checks if the provided password matches the stored hash.
func (u *User) VerifyPassword(password string) bool {
	if u.PasswordHash == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password))
	return err == nil
}
