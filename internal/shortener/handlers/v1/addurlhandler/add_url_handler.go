// Package addurlhandler Добавление URL V1
package addurlhandler

import (
	"errors"
	"io"
	"net/http"
	neturl "net/url"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/repositoryerror"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
)

// AddHandler Обработчик добавления нового URL
type AddHandler struct {
	urlRepository     url.URLRepositoryInterface
	userService       user.UserServiceInterface
	shortURLGenerator urlgen.ShortURLGeneratorInterface
	logger            *zap.Logger
	baseURL           string
}

// NewAddHandler missing godoc.
func NewAddHandler(
	urlRepository url.URLRepositoryInterface,
	userService user.UserServiceInterface,
	shortURLGenerator urlgen.ShortURLGeneratorInterface,
	zapLogger *zap.Logger,
	baseURL string,
) *AddHandler {
	return &AddHandler{
		urlRepository:     urlRepository,
		userService:       userService,
		shortURLGenerator: shortURLGenerator,
		logger:            zapLogger,
		baseURL:           baseURL,
	}
}

// AddURL Добавляет новый URL.
// Алгоритм работы:
// Генерирует для URL его короткую версию
// Сохраняет в базу URL
// Прикрепляет сохраненный URL к пользователю
//
// На вход приходит строка URL
// https://www.google1.ru/
// На выходе - готовая ссылка для редиректа
// http://localhost:8080/HnsSMA
func (handler *AddHandler) AddURL(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

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
	} else {
		if userID != "" {
			err = handler.userService.AddShortURLsToUser(userID, []string{shortKey})
			if err != nil {
				handler.logger.Error("add short url to user error", zap.String("error", err.Error()))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(status)

	_, err = w.Write([]byte(handler.baseURL + "/" + shortKey))

	if err != nil {
		handler.logger.Error("write response error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
