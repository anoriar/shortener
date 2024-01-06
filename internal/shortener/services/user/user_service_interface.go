// Package user Сервис работы с пользователем
package user

// UserServiceInterface Сервис работы с пользователем
//
//go:generate mockgen -source=user_service_interface.go -destination=mock/user_service.go -package=mock UserServiceInterface
type UserServiceInterface interface {
	// AddShortURLsToUser Добавляет URL к пользователю
	AddShortURLsToUser(userID string, shortURLs []string) error
	// GetUserShortURLs Получение всех URL, которые создал пользователь
	GetUserShortURLs(userID string) ([]string, error)
}
