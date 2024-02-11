// Package usecases Редирект на URL
package usecases

import (
	"strings"

	"github.com/anoriar/shortener/internal/shortener/domainerror"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/repository/url"
)

// GetURLService Обработчик редиректа по короткому URLу
type GetURLService struct {
	urlRepository url.URLRepositoryInterface
	logger        *zap.Logger
}

// NewGetURLService missing godoc.
func NewGetURLService(urlRepository url.URLRepositoryInterface, logger *zap.Logger) *GetURLService {
	return &GetURLService{
		urlRepository: urlRepository,
		logger:        logger,
	}
}

// GetURL получает URL из БД по короткому URL и осуществляет редирект по нему
//
// На вход в URLе приходит сокращенный URL: JRU9a8
// На выход: редирект с сокращенная версией URL
func (handler *GetURLService) GetURL(path string) (string, error) {

	shortKey := strings.Trim(path, "/")
	if shortKey == "" {
		return "", domainerror.ErrNotValidURL
	}

	url, err := handler.urlRepository.FindURLByShortURL(shortKey)
	if err != nil {
		handler.logger.Error("get url error", zap.String("error", err.Error()))
		return "", domainerror.ErrInternal
	}
	if url == nil {
		return "", domainerror.ErrURLNotExist
	}
	if url.IsDeleted {
		return "", domainerror.ErrURLDeleted
	}

	return url.OriginalURL, nil
}
