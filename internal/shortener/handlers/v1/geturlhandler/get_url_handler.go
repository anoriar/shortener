package geturlhandler

import (
	"github.com/anoriar/shortener/internal/shortener/repository"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type GetHandler struct {
	urlRepository repository.URLRepositoryInterface
	logger        *zap.Logger
}

func NewGetHandler(urlRepository repository.URLRepositoryInterface, logger *zap.Logger) *GetHandler {
	return &GetHandler{
		urlRepository: urlRepository,
		logger:        logger,
	}
}

func (handler *GetHandler) GetURL(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain")

	shortKey := strings.Trim(req.URL.Path, "/")
	if shortKey == "" {
		http.Error(w, "Short key is empty", http.StatusBadRequest)
		return
	}

	url, err := handler.urlRepository.FindURLByShortURL(shortKey)
	if err != nil {
		handler.logger.Error("get url error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if url == nil {
		http.Error(w, "URL does not exists", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
