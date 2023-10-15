package addurlhandler

import (
	"errors"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/repository/repositoryerror"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	neturl "net/url"
)

type AddHandler struct {
	urlRepository     repository.URLRepositoryInterface
	shortURLGenerator urlgen.ShortURLGeneratorInterface
	logger            *zap.Logger
	baseURL           string
}

func NewAddHandler(
	urlRepository repository.URLRepositoryInterface,
	shortURLGenerator urlgen.ShortURLGeneratorInterface,
	zapLogger *zap.Logger,
	baseURL string,
) *AddHandler {
	return &AddHandler{
		urlRepository:     urlRepository,
		shortURLGenerator: shortURLGenerator,
		logger:            zapLogger,
		baseURL:           baseURL,
	}
}

func (handler *AddHandler) AddURL(w http.ResponseWriter, req *http.Request) {
	status := http.StatusCreated

	url, err := io.ReadAll(req.Body)
	if err != nil {
		handler.logger.Error("read request body error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	parsedURL, err := neturl.Parse(string(url))
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "Not valid URL", http.StatusBadRequest)
		return
	}

	shortKey, err := handler.shortURLGenerator.GenerateShortURL()
	if err != nil {
		handler.logger.Error("generate URL error", zap.String("error", err.Error()))
		http.Error(w, "Not valid URL", http.StatusBadRequest)
		return
	}

	err = handler.urlRepository.AddURL(
		&entity.URL{
			UUID:        uuid.NewString(),
			ShortURL:    shortKey,
			OriginalURL: string(url),
		})

	if err != nil {
		if errors.Is(err, repositoryerror.ErrConflict) {
			existedURL, err := handler.urlRepository.FindURLByOriginalURL(req.Context(), string(url))
			if err != nil {
				handler.logger.Error("find existed URL error", zap.String("error", err.Error()))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			shortKey = existedURL.ShortURL
			status = http.StatusConflict
		} else {
			handler.logger.Error("add URL error", zap.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(status)
	w.Header().Set("content-type", "text/plain")

	_, err = w.Write([]byte(handler.baseURL + "/" + shortKey))
	if err != nil {
		handler.logger.Error("write response error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
