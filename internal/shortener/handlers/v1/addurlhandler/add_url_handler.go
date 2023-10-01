package addurlhandler

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/util"
	"github.com/google/uuid"
	"io"
	"net/http"
	neturl "net/url"
)

type AddHandler struct {
	urlRepository repository.URLRepositoryInterface
	keyGen        util.KeyGenInterface
	baseURL       string
}

func InitializeAddHandler(cnf *config.Config, repository repository.URLRepositoryInterface) *AddHandler {
	return NewAddHandler(repository, util.NewKeyGen(), cnf.BaseURL)
}

func NewAddHandler(urlRepository repository.URLRepositoryInterface, keyGen util.KeyGenInterface, baseURL string) *AddHandler {
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
	_, err = handler.urlRepository.AddURL(
		&entity.Url{
			Uuid:        uuid.NewString(),
			ShortURL:    shortKey,
			OriginalURL: string(url),
		})

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
