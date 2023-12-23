package getuserurlshandler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/repository"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/getuserurlshandler/internal/factory"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/services/user"
)

// GetUserURLsHandler missing godoc.
type GetUserURLsHandler struct {
	urlRepository   url.URLRepositoryInterface
	userService     user.UserServiceInterface
	responseFactory *factory.GetUSerURLsResponseFactory
	logger          *zap.Logger
}

// NewGetUserURLsHandler missing godoc.
func NewGetUserURLsHandler(
	urlRepository url.URLRepositoryInterface,
	userService user.UserServiceInterface,
	responseFactory *factory.GetUSerURLsResponseFactory,
	logger *zap.Logger,
) *GetUserURLsHandler {
	return &GetUserURLsHandler{
		urlRepository:   urlRepository,
		userService:     userService,
		responseFactory: responseFactory,
		logger:          logger,
	}
}

// GetUserURLs missing godoc.
func (handler *GetUserURLsHandler) GetUserURLs(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	w.Header().Set("content-type", "application/json")

	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	shortURLs, err := handler.userService.GetUserShortURLs(userID)
	if err != nil {
		handler.logger.Error("get short URLs from user error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(shortURLs) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resultURLs, err2 := handler.urlRepository.GetURLsByQuery(req.Context(), repository.Query{
		ShortURLs:    shortURLs,
		OriginalURLs: nil,
	})
	if err2 != nil {
		handler.logger.Error("get URLs error", zap.String("error", err2.Error()))
		http.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	if len(resultURLs) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response := handler.responseFactory.CreateResponse(resultURLs)
	jsonResult, err := json.Marshal(response)
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
