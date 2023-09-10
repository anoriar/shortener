package addurlhandler

import (
	"github.com/anoriar/shortener/internal/storage"
	"github.com/anoriar/shortener/internal/util"
	"io"
	"net/http"
	neturl "net/url"
)

type AddHandler struct {
	urlRepository storage.URLStorageInterface
	keyGen        util.KeyGenInterface
	baseURL       string
}

func NewAddHandler(urlRepository storage.URLStorageInterface, keyGen util.KeyGenInterface, baseURL string) *AddHandler {
	return &AddHandler{
		urlRepository: urlRepository,
		keyGen:        keyGen,
		baseURL:       baseURL,
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

	_, err = w.Write([]byte(handler.baseURL + "/" + shortKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
