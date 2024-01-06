// Package deleteuserurls Удаление URL, которые создал пользователь
package deleteuserurls

import (
	"context"
)

// DeleteUserURLsInterface Удаление URL, которые создал пользователь
//
//go:generate mockgen -source=delete_user_urls_interface.go -destination=mock/delete_user_urls_mock.go -package=mock DeleteUserURLsInterface
type DeleteUserURLsInterface interface {
	// DeleteUserURLs Удаляет URL, которые привязаны к пользователю. Под удалением подразумевается проставить флаг is_deleted в хранилище
	DeleteUserURLs(ctx context.Context, userID string, shortURLs []string) error
}
