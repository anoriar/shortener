// Package addurlbatchhander Добавление урлов пачкой
package addurlbatchhander

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/factory"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/validator"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/services/user"
)

// AddURLBatchHandler Обработчик добавления урлов пачкой
type AddURLBatchHandler struct {
	urlRepository              url.URLRepositoryInterface
	userService                user.UserServiceInterface
	addURLBatchFactory         *factory.AddURLEntityFactory
	addURLBatchResponseFactory *factory.AddURLBatchResponseFactory
	logger                     *zap.Logger
	validator                  *validator.AddURLBatchValidator
}

// NewAddURLBatchHandler missing godoc.
func NewAddURLBatchHandler(
	urlRepository url.URLRepositoryInterface,
	userService user.UserServiceInterface,
	addURLBatchFactory *factory.AddURLEntityFactory,
	addURLBatchResponseFactory *factory.AddURLBatchResponseFactory,
	logger *zap.Logger,
	validator *validator.AddURLBatchValidator,
) *AddURLBatchHandler {
	return &AddURLBatchHandler{
		urlRepository:              urlRepository,
		userService:                userService,
		addURLBatchFactory:         addURLBatchFactory,
		addURLBatchResponseFactory: addURLBatchResponseFactory,
		logger:                     logger,
		validator:                  validator,
	}
}

// AddURLBatch добавляет несколько URL на основе входящего запроса.
//
// Процесс работы функции включает следующие шаги:
// 1. Генерация короткой версии для каждого URL.
// 2. Сохранение всех URL в базу данных.
// 3. Прикрепление сохранённых URL к конкретному пользователю.
// 4. Сопоставление входных и выходных данных по correlation_id и возврат сгенерированных коротких ссылок.
//
// Формат входных данных:
// [
//
//	{
//	  "correlation_id": "by4564trg",
//	  "original_url": "https://practicum3.yandex.ru"
//	},
//	...
//
// ]
//
// Формат выходных данных:
// [
//
//	{
//	  "correlation_id": "by4564trg",
//	  "short_url": "http://localhost:8080/Ytq3tY"
//	},
//	...
//
// ]
//
// Параметр correlation_id используется для сопоставления входных и выходных URL.
// Обратите внимание, что это поле не используется в базе данных.
func (handler *AddURLBatchHandler) AddURLBatch(w http.ResponseWriter, req *http.Request) {
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

	var requestItems []request.AddURLBatchRequestDTO

	err = json.Unmarshal(requestBody, &requestItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = handler.validator.Validate(requestItems); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlsMap := handler.addURLBatchFactory.CreateURLsFromBatchRequest(requestItems)
	urls := make([]entity.URL, 0, len(urlsMap))
	for _, urlEntity := range urlsMap {
		urls = append(urls, urlEntity)
	}

	err = handler.urlRepository.AddURLBatch(req.Context(), urls)
	if err != nil {
		handler.logger.Error("batch add error", zap.String("error", err.Error()))
		http.Error(w, "batch add error", http.StatusBadRequest)
		return
	} else {
		if userID != "" {
			shortKeys := make([]string, 0, len(urlsMap))
			for _, val := range urlsMap {
				shortKeys = append(shortKeys, val.ShortURL)
			}

			err = handler.userService.AddShortURLsToUser(userID, shortKeys)
			if err != nil {
				handler.logger.Error("add short url to user error", zap.String("error", err.Error()))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	response := handler.addURLBatchResponseFactory.CreateResponse(urlsMap, requestItems)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

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
}
