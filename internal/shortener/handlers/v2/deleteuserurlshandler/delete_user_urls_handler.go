// Package deleteuserurlshandler Модуль удаления URL, которые создал пользователь
package deleteuserurlshandler

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/context"
	deleteurlsprocessor "github.com/anoriar/shortener/internal/shortener/processors/deleteuserurlsprocessor"
	"github.com/anoriar/shortener/internal/shortener/processors/deleteuserurlsprocessor/message"
)

// DeleteUserURLsHandler Обработчик удаления URL, которые создал пользователь
type DeleteUserURLsHandler struct {
	deleteUserURLsProcessor *deleteurlsprocessor.DeleteUserURLsProcessor
	logger                  *zap.Logger
}

// NewDeleteUserURLsHandler missing godoc.
func NewDeleteUserURLsHandler(deleteUserURLsProcessor *deleteurlsprocessor.DeleteUserURLsProcessor, logger *zap.Logger) *DeleteUserURLsHandler {
	return &DeleteUserURLsHandler{deleteUserURLsProcessor: deleteUserURLsProcessor, logger: logger}
}

// DeleteUserURLs Удаляет несколько URL, которые создал пользователь
// Под удалением подразумевается - пометить в БД фла is_deleted=true
// Обработка происходит асинхронно
// На вход
// [
//
//	"6qxTVvsy",
//	"RTfd56hn",
//	"Jlfd67ds"
//
// ]
func (handler *DeleteUserURLsHandler) DeleteUserURLs(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		handler.logger.Error("read request error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var shortURLs []string

	err = json.Unmarshal(requestBody, &shortURLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if userID != "" && len(shortURLs) > 0 {
		handler.deleteUserURLsProcessor.AddMessage(message.DeleteUserURLsMessage{
			UserID:    userID,
			ShortURLs: shortURLs,
		})
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
}
