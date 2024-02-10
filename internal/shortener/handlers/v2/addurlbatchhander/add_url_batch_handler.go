// Package addurlbatchhander Добавление урлов пачкой
package addurlbatchhander

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/anoriar/shortener/internal/shortener/usecases/addurlbatch"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/dto/request"
)

// AddURLBatchHandler Обработчик добавления урлов пачкой
type AddURLBatchHandler struct {
	logger             *zap.Logger
	addURLBatchService *addurlbatch.AddURLBatchService
}

// NewAddURLBatchHandler missing godoc.
func NewAddURLBatchHandler(
	logger *zap.Logger,
	addURLBatchService *addurlbatch.AddURLBatchService,
) *AddURLBatchHandler {
	return &AddURLBatchHandler{
		logger:             logger,
		addURLBatchService: addURLBatchService,
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

	response, err := handler.addURLBatchService.AddURLBatch(req.Context(), requestItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
