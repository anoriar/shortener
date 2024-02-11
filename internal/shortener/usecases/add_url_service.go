package usecases

import (
	"context"
	"errors"
	neturl "net/url"

	"github.com/google/uuid"
	"go.uber.org/zap"

	internalctx "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/domainerror"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/repositoryerror"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
)

// AddURLService Обработчик добавления нового URL
type AddURLService struct {
	urlRepository     url.URLRepositoryInterface
	shortURLGenerator urlgen.ShortURLGeneratorInterface
	userService       user.UserServiceInterface
	logger            *zap.Logger
	baseURL           string
}

// NewAddURLService missing godoc.
func NewAddURLService(
	urlRepository url.URLRepositoryInterface,
	shortURLGenerator urlgen.ShortURLGeneratorInterface,
	userService user.UserServiceInterface,
	logger *zap.Logger,
	baseURL string,
) *AddURLService {
	return &AddURLService{
		urlRepository:     urlRepository,
		shortURLGenerator: shortURLGenerator,
		userService:       userService,
		logger:            logger,
		baseURL:           baseURL,
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
func (service AddURLService) AddURL(ctx context.Context, addURLRequestDto request.AddURLRequestDto) (*response.AddURLResponseDto, error) {
	userID := ""
	userIDCtxParam := ctx.Value(internalctx.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	parsedURL, err := neturl.Parse(addURLRequestDto.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, domainerror.ErrNotValidURL
	}

	shortKey, err := service.shortURLGenerator.GenerateShortURL()
	if err != nil {
		service.logger.Error("short URL generation error", zap.String("error", err.Error()))
		return nil, domainerror.ErrInternal
	}

	err = service.urlRepository.AddURL(&entity.URL{
		UUID:        uuid.NewString(),
		ShortURL:    shortKey,
		OriginalURL: addURLRequestDto.URL,
	})

	if err != nil {
		if errors.Is(err, repositoryerror.ErrConflict) {
			existedURL, err := service.urlRepository.FindURLByOriginalURL(ctx, addURLRequestDto.URL)
			if err != nil {
				service.logger.Error("find existed URL error", zap.String("error", err.Error()))
				return nil, domainerror.ErrURLExists
			}
			shortKey = existedURL.ShortURL
			return &response.AddURLResponseDto{
				Result: service.baseURL + "/" + shortKey,
			}, domainerror.ErrURLExists
		} else {
			service.logger.Error("add URL error", zap.String("error", err.Error()))
			return nil, domainerror.ErrInternal
		}
	} else {
		if userID != "" {
			err = service.userService.AddShortURLsToUser(userID, []string{shortKey})
			if err != nil {
				service.logger.Error("add short url to user error", zap.String("error", err.Error()))
				return nil, domainerror.ErrInternal
			}
		}
	}
	responseDTO := response.AddURLResponseDto{
		Result: service.baseURL + "/" + shortKey,
	}

	return &responseDTO, nil
}
