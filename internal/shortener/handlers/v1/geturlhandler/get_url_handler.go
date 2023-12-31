// Package geturlhandler Редирект на URL
package geturlhandler

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/repository/url"
)

// GetHandler Обработчик редиректа по короткому URLу
type GetHandler struct {
	urlRepository url.URLRepositoryInterface
	logger        *zap.Logger
}

// NewGetHandler missing godoc.
func NewGetHandler(urlRepository url.URLRepositoryInterface, logger *zap.Logger) *GetHandler {
	return &GetHandler{
		urlRepository: urlRepository,
		logger:        logger,
	}
}

// GetURL получает URL из БД по короткому URL и осуществляет редирект по нему
//
// На вход в URLе приходит сокращенный URL: JRU9a8
// На выход: 307 редирект с сокращенной версией URL в заголовке Location
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
	if url.IsDeleted {
		http.Error(w, "URL deleted", http.StatusGone)
		return
	}

	w.Header().Set("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
