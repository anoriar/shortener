package user

import "github.com/anoriar/shortener/internal/shortener/entity"

//go:generate mockgen -source=user_repository_interface.go -destination=mock/user_repository.go -package=mock UserRepositoryInterface
type UserRepositoryInterface interface {
	AddUser(user entity.User) error
	UpdateUser(user entity.User) error
	FindUserByID(userID string) (entity.User, bool, error)
}
