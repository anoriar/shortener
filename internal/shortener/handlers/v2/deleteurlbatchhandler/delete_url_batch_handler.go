package deleteurlbatchhandler

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type DeleteURLBatchHandler struct {
	urlRepository url.URLRepositoryInterface
	logger        *zap.Logger
}

func NewDeleteURLBatchHandler(urlRepository url.URLRepositoryInterface, logger *zap.Logger) *DeleteURLBatchHandler {
	return &DeleteURLBatchHandler{urlRepository: urlRepository, logger: logger}
}

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
		handler.logger.Error("batch add error", zap.String("error", err.Error()))
		http.Error(w, "batch add error", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
}
