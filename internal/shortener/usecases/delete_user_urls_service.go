// Package DeleteUserURLsService Модуль удаления URL, которые создал пользователь
package usecases

import (
	"context"

	"go.uber.org/zap"

	internalctx "github.com/anoriar/shortener/internal/shortener/context"
	deleteurlsprocessor "github.com/anoriar/shortener/internal/shortener/processors/deleteuserurlsprocessor"
	"github.com/anoriar/shortener/internal/shortener/processors/deleteuserurlsprocessor/message"
)

// DeleteUserURLsService Обработчик удаления URL, которые создал пользователь
type DeleteUserURLsService struct {
	deleteUserURLsProcessor *deleteurlsprocessor.DeleteUserURLsProcessor
	logger                  *zap.Logger
}

// NewDeleteUserURLsService missing godoc.
func NewDeleteUserURLsService(deleteUserURLsProcessor *deleteurlsprocessor.DeleteUserURLsProcessor, logger *zap.Logger) *DeleteUserURLsService {
	return &DeleteUserURLsService{deleteUserURLsProcessor: deleteUserURLsProcessor, logger: logger}
}

// DeleteUserURLs удаляет несколько URL, которые создал пользователь.
// Под удалением подразумевается пометка в БД флага is_deleted=true.
// Обработка происходит асинхронно.
//
// На вход принимает массив сокращенных версий URL, созданных пользователем:
// [
//
//	"6qxTVvsy",
//	"RTfd56hn",
//	"Jlfd67ds"
//
// ]
func (handler *DeleteUserURLsService) DeleteUserURLs(ctx context.Context, shortURLs []string) error {
	userID := ""
	userIDCtxParam := ctx.Value(internalctx.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID != "" && len(shortURLs) > 0 {
		handler.deleteUserURLsProcessor.AddMessage(message.DeleteUserURLsMessage{
			UserID:    userID,
			ShortURLs: shortURLs,
		})
	}

	return nil
}
