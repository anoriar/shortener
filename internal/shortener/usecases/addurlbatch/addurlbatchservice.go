// Package addurlbatch Добавление урлов пачкой
package addurlbatch

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/domainerror"
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/util"

	internalctx "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/usecases/addurlbatch/internal/factory"
	"github.com/anoriar/shortener/internal/shortener/usecases/addurlbatch/internal/validator"
)

// AddURLBatchService Обработчик добавления урлов пачкой
type AddURLBatchService struct {
	urlRepository              url.URLRepositoryInterface
	userService                user.UserServiceInterface
	addURLBatchFactory         *factory.AddURLEntityFactory
	addURLBatchResponseFactory *factory.AddURLBatchResponseFactory
	logger                     *zap.Logger
	validator                  *validator.AddURLBatchValidator
}

// NewAddURLBatchService missing godoc.
func NewAddURLBatchService(
	urlRepository url.URLRepositoryInterface,
	userService user.UserServiceInterface,
	keyGen util.KeyGenInterface,
	baseURL string,
	logger *zap.Logger,
) *AddURLBatchService {
	return &AddURLBatchService{
		urlRepository:              urlRepository,
		userService:                userService,
		addURLBatchFactory:         factory.NewAddURLBatchFactory(keyGen),
		addURLBatchResponseFactory: factory.NewAddURLBatchResponseFactory(baseURL),
		logger:                     logger,
		validator:                  validator.NewAddURLBatchValidator(),
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
func (service *AddURLBatchService) AddURLBatch(ctx context.Context, requestItems []request.AddURLBatchRequestDTO) ([]response.AddURLBatchResponseDTO, error) {
	userID := ""
	userIDCtxParam := ctx.Value(internalctx.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if err := service.validator.Validate(requestItems); err != nil {
		return nil, fmt.Errorf("%w: %s", domainerror.ErrNotValidRequest, err.Error())
	}

	urlsMap := service.addURLBatchFactory.CreateURLsFromBatchRequest(requestItems)
	urls := make([]entity.URL, 0, len(urlsMap))
	for _, urlEntity := range urlsMap {
		urls = append(urls, urlEntity)
	}

	err := service.urlRepository.AddURLBatch(ctx, urls)
	if err != nil {
		service.logger.Error("batch add error", zap.String("error", err.Error()))
		return nil, domainerror.ErrInternal
	} else {
		if userID != "" {
			shortKeys := make([]string, 0, len(urlsMap))
			for _, val := range urlsMap {
				shortKeys = append(shortKeys, val.ShortURL)
			}

			err = service.userService.AddShortURLsToUser(userID, shortKeys)
			if err != nil {
				service.logger.Error("add short url to user error", zap.String("error", err.Error()))
				return nil, domainerror.ErrInternal
			}
		}
	}

	responseItems := service.addURLBatchResponseFactory.CreateResponse(urlsMap, requestItems)

	return responseItems, nil
}
