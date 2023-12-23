package inmemory

import (
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/repositoryerror"
)

// InMemoryUserRepository missing godoc.
type InMemoryUserRepository struct {
	users map[string]entity.User
}

// NewInMemoryUserRepository missing godoc.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]entity.User)}
}

// AddUser missing godoc.
func (repository InMemoryUserRepository) AddUser(user entity.User) error {
	if _, exists := repository.users[user.UUID]; exists {
		return repositoryerror.ErrConflict
	}
	repository.users[user.UUID] = user

	return nil
}

// UpdateUser missing godoc.
func (repository InMemoryUserRepository) UpdateUser(user entity.User) error {
	if _, exists := repository.users[user.UUID]; !exists {
		return repositoryerror.ErrNotFound
	}
	repository.users[user.UUID] = user

	return nil
}

// FindUserByID missing godoc.
func (repository InMemoryUserRepository) FindUserByID(userID string) (entity.User, bool, error) {
	user, exists := repository.users[userID]
	if !exists {
		return entity.User{}, false, nil
	}
	return user, true, nil
}
