package deleteuserurlshandler

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/repository/user"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type DeleteUserURLsHandler struct {
	urlRepository  url.URLRepositoryInterface
	userRepository user.UserRepositoryInterface
	logger         *zap.Logger
}

func NewDeleteUserURLsHandler(urlRepository url.URLRepositoryInterface, userRepository user.UserRepositoryInterface, logger *zap.Logger) *DeleteUserURLsHandler {
	return &DeleteUserURLsHandler{urlRepository: urlRepository, userRepository: userRepository, logger: logger}
}

func (handler *DeleteUserURLsHandler) DeleteUserURLs(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		handler.logger.Error("read request error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var shortURLs []string

	err = json.Unmarshal(requestBody, &shortURLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, exist, err := handler.userRepository.FindUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if exist && len(user.SavedURLIDs) > 0 {
		var shortURLsForDelete []string

		for _, requestShortURL := range shortURLs {
			if _, exists := user.SavedURLIDs.FindShortURL(requestShortURL); exists {
				shortURLsForDelete = append(shortURLsForDelete, requestShortURL)
			}
		}
		if len(shortURLsForDelete) > 0 {
			err = handler.urlRepository.UpdateIsDeletedBatch(req.Context(), shortURLsForDelete, true)
			if err != nil {
				handler.logger.Error("batch delete error", zap.String("error", err.Error()))
				http.Error(w, "batch delete error", http.StatusBadRequest)
				return
			}
		}
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
}
