// Package geturlhandler Редирект на URL
package geturlhandler

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/domainerror"
	"github.com/anoriar/shortener/internal/shortener/usecases"
)

// GetHandler Обработчик редиректа по короткому URLу
type GetHandler struct {
	logger        *zap.Logger
	getURLService *usecases.GetURLService
}

// NewGetHandler missing godoc.
func NewGetHandler(logger *zap.Logger, getURLService *usecases.GetURLService) *GetHandler {
	return &GetHandler{
		logger:        logger,
		getURLService: getURLService,
	}
}

// GetURL получает URL из БД по короткому URL и осуществляет редирект по нему
//
// На вход в URLе приходит сокращенный URL: JRU9a8
// На выход: 307 редирект с сокращенной версией URL в заголовке Location
func (handler *GetHandler) GetURL(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain")

	url, err := handler.getURLService.GetURL(req.URL.Path)
	if err != nil {
		switch {
		case errors.Is(err, domainerror.ErrURLDeleted):
			http.Error(w, err.Error(), http.StatusGone)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
