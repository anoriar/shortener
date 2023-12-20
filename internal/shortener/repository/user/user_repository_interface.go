// Package user Пакет работы с хранилищем пользователей
package user

import "github.com/anoriar/shortener/internal/shortener/entity"

// UserRepositoryInterface Интерфейс работы с хранилищем пользователей
//
//go:generate mockgen -source=user_repository_interface.go -destination=mock/user_repository.go -package=mock UserRepositoryInterface
type UserRepositoryInterface interface {
	// AddUser Добавление пользователя в хранилище
	AddUser(user entity.User) error
	// UpdateUser Обновление информации о пользователе в хранилище
	UpdateUser(user entity.User) error
	// FindUserByID Получение пользователя по ID
	FindUserByID(userID string) (entity.User, bool, error)
}
