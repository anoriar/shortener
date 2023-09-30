package geturlhandler

import (
	"github.com/anoriar/shortener/internal/shortener/repository"
	"net/http"
	"strings"
)

type GetHandler struct {
	urlRepository repository.URLRepositoryInterface
}

func InitializeGetHandler(repository repository.URLRepositoryInterface) *GetHandler {
	return NewGetHandler(repository)
}

func NewGetHandler(urlRepository repository.URLRepositoryInterface) *GetHandler {
	return &GetHandler{
		urlRepository: urlRepository,
	}
}

func (handler *GetHandler) GetURL(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain")

	shortKey := strings.Trim(req.URL.Path, "/")
	if shortKey == "" {
		http.Error(w, "Short key is empty", http.StatusBadRequest)
		return
	}

	url, err := handler.urlRepository.FindURLByKey(shortKey)
	if err != nil {
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
