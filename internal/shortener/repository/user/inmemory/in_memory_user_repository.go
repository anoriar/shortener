package inmemory

import (
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/repositoryerror"
)

type InMemoryUserRepository struct {
	users map[string]entity.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]entity.User)}
}

func (repository InMemoryUserRepository) AddUser(user entity.User) error {
	if _, exists := repository.users[user.UUID]; exists {
		return repositoryerror.ErrConflict
	}
	repository.users[user.UUID] = user

	return nil
}

func (repository InMemoryUserRepository) UpdateUser(user entity.User) error {
	if _, exists := repository.users[user.UUID]; !exists {
		return repositoryerror.ErrNotFound
	}
	repository.users[user.UUID] = user

	return nil
}

func (repository InMemoryUserRepository) FindUserByID(userID string) (entity.User, bool, error) {
	user, exists := repository.users[userID]
	if !exists {
		return entity.User{}, false, nil
	}
	return user, true, nil
}
