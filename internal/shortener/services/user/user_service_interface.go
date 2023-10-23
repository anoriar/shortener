package user

//go:generate mockgen -source=user_service_interface.go -destination=mock/user_service.go -package=mock UserServiceInterface
type UserServiceInterface interface {
	AddShortURLsToUser(userID string, shortURLs []string) error
	GetUserShortURLs(userID string) ([]string, error)
}
