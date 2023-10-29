package deleteuserurls

import (
	"context"
)

//go:generate mockgen -source=delete_user_urls_interface.go -destination=mock/delete_user_urls_mock.go -package=mock DeleteUserURLsInterface
type DeleteUserURLsInterface interface {
	DeleteUserURLs(ctx context.Context, userID string, shortURLs []string) error
}
