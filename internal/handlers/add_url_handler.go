package handlers

import (
	"github.com/anoriar/shortener/internal/storage"
	"github.com/anoriar/shortener/internal/util"
	"io"
	"net/http"
	neturl "net/url"
)

type AddHandler struct {
	urlRepository storage.URLRepositoryInterface
	keyGen        util.KeyGenInterface
}

func NewAddHandler(urlRepository storage.URLRepositoryInterface, keyGen util.KeyGenInterface) *AddHandler {
	return &AddHandler{
		urlRepository: urlRepository,
		keyGen:        keyGen,
	}
}

func (handler *AddHandler) AddURL(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain")

	url, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	parsedURL, err := neturl.Parse(string(url))
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "Not valid URL", http.StatusBadRequest)
		return
	}

	shortKey := handler.keyGen.Generate()
	err = handler.urlRepository.AddURL(string(url), shortKey)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte("http://localhost:8080/" + shortKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
