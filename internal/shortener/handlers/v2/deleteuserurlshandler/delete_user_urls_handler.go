// Package deleteuserurlshandler Модуль удаления URL, которые создал пользователь
package deleteuserurlshandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/anoriar/shortener/internal/shortener/usecases"

	"go.uber.org/zap"
)

// DeleteUserURLsHandler Обработчик удаления URL, которые создал пользователь
type DeleteUserURLsHandler struct {
	logger                *zap.Logger
	deleteUserURLsService *usecases.DeleteUserURLsService
}

// NewDeleteUserURLsHandler missing godoc.
func NewDeleteUserURLsHandler(deleteUserURLsService *usecases.DeleteUserURLsService, logger *zap.Logger) *DeleteUserURLsHandler {
	return &DeleteUserURLsHandler{deleteUserURLsService: deleteUserURLsService, logger: logger}
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
func (handler *DeleteUserURLsHandler) DeleteUserURLs(w http.ResponseWriter, req *http.Request) {
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
	err = handler.deleteUserURLsService.DeleteUserURLs(req.Context(), shortURLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
}
