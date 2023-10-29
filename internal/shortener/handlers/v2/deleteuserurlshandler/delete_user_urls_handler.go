package deleteuserurlshandler

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/context"
	deleteurlsprocessor "github.com/anoriar/shortener/internal/shortener/processors/deleteuserurlsprocessor"
	"github.com/anoriar/shortener/internal/shortener/processors/deleteuserurlsprocessor/message"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type DeleteUserURLsHandler struct {
	deleteUserURLsProcessor *deleteurlsprocessor.DeleteUserURLsProcessor
	logger                  *zap.Logger
}

func NewDeleteUserURLsHandler(deleteUserURLsProcessor *deleteurlsprocessor.DeleteUserURLsProcessor, logger *zap.Logger) *DeleteUserURLsHandler {
	return &DeleteUserURLsHandler{deleteUserURLsProcessor: deleteUserURLsProcessor, logger: logger}
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
	if userID != "" && len(shortURLs) > 0 {
		//#MENTOR: Не понял в задании
		// "Для максимального наполнения буфера объектов обновления используйте паттерн fanIn"
		// С точки зрения обще проектировки: сервер создает горутины - (fan-out)
		// Здесь засовываем все сообщения в один канал - а обрабатывается он асинхронно
		// Может в задании что-то другое подразумевалось?
		handler.deleteUserURLsProcessor.AddMessage(message.DeleteUserURLsMessage{
			UserID:    userID,
			ShortURLs: shortURLs,
		})
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
}
