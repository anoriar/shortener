package getuserurlbatch

import (
	"context"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/domainerror"
	"github.com/anoriar/shortener/internal/shortener/dto/response"

	internalctx "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/repository"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/usecases/getuserurlbatch/internal/factory"
)

// GetUserURLsService missing godoc.
type GetUserURLsService struct {
	urlRepository   url.URLRepositoryInterface
	userService     user.UserServiceInterface
	responseFactory *factory.GetUSerURLsResponseFactory
	logger          *zap.Logger
}

// NewGetUserURLsService missing godoc.
func NewGetUserURLsService(
	urlRepository url.URLRepositoryInterface,
	userService user.UserServiceInterface,
	logger *zap.Logger,
	baseURL string,
) *GetUserURLsService {
	return &GetUserURLsService{
		urlRepository:   urlRepository,
		userService:     userService,
		responseFactory: factory.NewGetUSerURLsResponseFactory(baseURL),
		logger:          logger,
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
func (service *GetUserURLsService) GetUserURLs(ctx context.Context) ([]response.GetUserURLsResponseDTO, error) {
	userID := ""
	userIDCtxParam := ctx.Value(internalctx.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID == "" {
		return nil, domainerror.ErrUserUnauthorized
	}

	shortURLs, err := service.userService.GetUserShortURLs(userID)
	if err != nil {
		service.logger.Error("get short URLs from user error", zap.String("error", err.Error()))
		return nil, domainerror.ErrInternal
	}

	if len(shortURLs) == 0 {
		return nil, domainerror.ErrUserUnauthorized
	}

	resultURLs, err2 := service.urlRepository.GetURLsByQuery(ctx, repository.Query{
		ShortURLs:    shortURLs,
		OriginalURLs: nil,
	})
	if err2 != nil {
		service.logger.Error("get URLs error", zap.String("error", err2.Error()))
		return nil, domainerror.ErrInternal
	}

	if len(resultURLs) == 0 {
		return nil, domainerror.ErrUserUnauthorized
	}

	responseItems := service.responseFactory.CreateResponse(resultURLs)

	return responseItems, nil
}
