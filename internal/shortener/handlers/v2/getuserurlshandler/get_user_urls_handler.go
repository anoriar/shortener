package getuserurlshandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/anoriar/shortener/internal/shortener/domainerror"
	"github.com/anoriar/shortener/internal/shortener/usecases/getuserurlbatch"

	"go.uber.org/zap"
)

// GetUserURLsHandler missing godoc.
type GetUserURLsHandler struct {
	logger             *zap.Logger
	getUserURLsService *getuserurlbatch.GetUserURLsService
}

// NewGetUserURLsHandler missing godoc.
func NewGetUserURLsHandler(
	logger *zap.Logger,
	getUserURLsService *getuserurlbatch.GetUserURLsService,
) *GetUserURLsHandler {
	return &GetUserURLsHandler{
		logger:             logger,
		getUserURLsService: getUserURLsService,
	}
}

// GetUserURLs получает URL, которые создал пользователь
// Формат выходных данных:
// [
//
//	{
//	  "original_url": "https://www.google1.ru/",
//	  "short_url": "http://localhost:8080/Ytq3tY"
//	},
//	...
//
// ]
func (handler *GetUserURLsHandler) GetUserURLs(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	responseItems, err := handler.getUserURLsService.GetUserURLs(req.Context())
	if err != nil {
		switch {
		case errors.Is(err, domainerror.ErrUserUnauthorized):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}

	jsonResult, err := json.Marshal(responseItems)
	if err != nil {
		handler.logger.Error("marshal error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(jsonResult)
	if err != nil {
		handler.logger.Error("write response error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
