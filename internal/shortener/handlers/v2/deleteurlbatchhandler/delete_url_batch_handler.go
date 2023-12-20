// Package deleteurlbatchhandler Удаление URL пачкой
package deleteurlbatchhandler

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/repository/url"
)

// DeleteURLBatchHandler Обработчик удаления URL пачкой
type DeleteURLBatchHandler struct {
	urlRepository url.URLRepositoryInterface
	logger        *zap.Logger
}

func NewDeleteURLBatchHandler(urlRepository url.URLRepositoryInterface, logger *zap.Logger) *DeleteURLBatchHandler {
	return &DeleteURLBatchHandler{urlRepository: urlRepository, logger: logger}
}

// DeleteURLBatch Удаляет несколько URL
// На вход: сокращенные версии URL
// [
//
//	"g95W3D",
//	"L7ibuA",
//	"TnZWBr"
//
// ]
func (handler *DeleteURLBatchHandler) DeleteURLBatch(w http.ResponseWriter, req *http.Request) {
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

	err = handler.urlRepository.DeleteURLBatch(req.Context(), shortURLs)
	if err != nil {
		handler.logger.Error("batch delete error", zap.String("error", err.Error()))
		http.Error(w, "batch delete error", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
}
