// Package addurlhandler Добавление URL V2
package addurlhandler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/domainerror"
	"github.com/anoriar/shortener/internal/shortener/usecases"

	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
)

// AddHandler Обработчик добавления нового URL
type AddHandler struct {
	urlRepository     url.URLRepositoryInterface
	shortURLGenerator urlgen.ShortURLGeneratorInterface
	userService       user.UserServiceInterface
	logger            *zap.Logger
	baseURL           string
	addURLService     *usecases.AddURLService
}

// NewAddHandler missing godoc.
func NewAddHandler(
	logger *zap.Logger,
	addURLService *usecases.AddURLService,
) *AddHandler {
	return &AddHandler{
		logger:        logger,
		addURLService: addURLService,
	}
}

// AddURL Добавляет новый URL.
// Алгоритм работы:
// - Генерирует для URL его короткую версию.
// - Сохраняет в базу URL.
// - Прикрепляет сохраненный URL к пользователю.
//
// На вход приходит:
//
//	{
//	   "url": "https://www.google1.ru/"
//	}
//
// На выходе - готовая ссылка для редиректа:
//
//	{
//	   "result": "http://localhost:8080/HnsSMA"
//	}
func (handler AddHandler) AddURL(w http.ResponseWriter, req *http.Request) {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		handler.logger.Error("read request error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addURLRequestDto := &request.AddURLRequestDto{}

	if err = json.Unmarshal(requestBody, addURLRequestDto); err != nil {
		handler.logger.Error("request unmarshal error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseDto, err := handler.addURLService.AddURL(req.Context(), *addURLRequestDto)
	if err != nil && responseDto == nil {
		switch {
		case errors.Is(err, domainerror.ErrNotValidURL):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, domainerror.ErrURLExists):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal server error", http.StatusBadRequest)
			return
		}
	}

	status := http.StatusCreated
	if err != nil && responseDto != nil {
		if errors.Is(err, domainerror.ErrURLExists) {
			status = http.StatusConflict
		}
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	jsonResult, err := json.Marshal(responseDto)
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
