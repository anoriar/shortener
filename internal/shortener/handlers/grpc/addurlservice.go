package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/repository/url"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	pb "github.com/anoriar/shortener/proto/generated/shortener/proto"
)

// AddURLService missing godoc.
type AddURLService struct {
	pb.UnimplementedAddURLServiceServer

	urlRepository     url.URLRepositoryInterface
	shortURLGenerator urlgen.ShortURLGeneratorInterface
	userService       user.UserServiceInterface
	logger            *zap.Logger
	baseURL           string
}

// NewAddURLService missing godoc.
// AddURLService missing godoc.
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
func (service AddURLService) AddURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	var response *pb.AddURLResponse

	return response, nil
}
