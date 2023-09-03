package handlers

import (
	"github.com/anoriar/shortener/internal/storage"
	"net/http"
	"strings"
)

type GetHandler struct {
	urlRepository storage.URLRepositoryInterface
}

func NewGetHandler(urlRepository storage.URLRepositoryInterface) *GetHandler {
	return &GetHandler{
		urlRepository: urlRepository,
	}
}

func (handler *GetHandler) GetURL(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain")

	shortKey := strings.Trim(req.URL.Path, "/")
	if shortKey == "" {
		http.Error(w, "Short key is empty", http.StatusBadRequest)
	}

	url, exists := handler.urlRepository.FindURLByKey(shortKey)
	if !exists {
		http.Error(w, "URL does not exists", http.StatusBadRequest)
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
