package repositories

import (
	"sync"

	"github.com/danigrb.dev/user-service/internal/database/interfaces"
)

var (
	userRepositoryInstance interfaces.UserRepository
	userRepositoryOnce     sync.Once
)

// Factory provides a centralized way to get repository instances
// This pattern allows for easy testing with mocks and centralized repository management
type Factory struct{}

// NewFactory creates a new repository factory
func NewFactory() *Factory {
	return &Factory{}
}

// GetUserRepository returns a UserRepository instance
// This uses the singleton pattern to ensure only one instance is created
func (f *Factory) GetUserRepository() interfaces.UserRepository {
	userRepositoryOnce.Do(func() {
		userRepositoryInstance = NewUserRepository()
	})
	return userRepositoryInstance
}

// SetUserRepository allows setting a custom UserRepository implementation
// This is particularly useful for testing with mocks
func (f *Factory) SetUserRepository(repo interfaces.UserRepository) {
	userRepositoryInstance = repo
	// Reset the once so GetUserRepository will use the new instance
	userRepositoryOnce = sync.Once{}
}
