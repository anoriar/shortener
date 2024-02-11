package grpc

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anoriar/shortener/internal/app"
	"github.com/anoriar/shortener/internal/shortener/domainerror"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/usecases"
	"github.com/anoriar/shortener/internal/shortener/usecases/addurlbatch"
	"github.com/anoriar/shortener/internal/shortener/usecases/getuserurlbatch"
	pb "github.com/anoriar/shortener/proto/generated/shortener/proto"
)

// URLServiceServer missing godoc.
type URLServiceServer struct {
	pb.UnimplementedURLServiceServer

	addURLService          *usecases.AddURLService
	addURLBatchService     *addurlbatch.AddURLBatchService
	getURLService          *usecases.GetURLService
	getUserURLBatchService *getuserurlbatch.GetUserURLsService
	deleteUserURLsService  *usecases.DeleteUserURLsService
	urlRepository          url.URLRepositoryInterface
	logger                 *zap.Logger
}

// NewURLServiceServer missing godoc.
// URLServiceServer missing godoc.
func NewURLServiceServer(
	app app.App,
) *URLServiceServer {
	return &URLServiceServer{
		addURLService:          app.AddURLServiceUC,
		addURLBatchService:     app.AddURLBatchServiceUC,
		getURLService:          app.GetURLServiceUC,
		getUserURLBatchService: app.GetUserURLsServiceUC,
		deleteUserURLsService:  app.DeleteUserURLsServiceUC,
		urlRepository:          app.URLRepository,
		logger:                 app.Logger,
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
func (service URLServiceServer) AddURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	var response pb.AddURLResponse

	urlDto, err := service.addURLService.AddURL(ctx, request.AddURLRequestDto{URL: in.GetUrl()})

	if err != nil && urlDto == nil {
		switch {
		case errors.Is(err, domainerror.ErrNotValidURL):
			return nil, status.Errorf(codes.InvalidArgument, `not valid URL`)
		case errors.Is(err, domainerror.ErrURLExists):
			return nil, status.Errorf(codes.AlreadyExists, `resource exists`)
		default:
			return nil, status.Errorf(codes.Internal, `internal error`)
		}
	}

	response.Result = urlDto.Result

	return &response, nil
}

// AddURLBatch добавляет несколько URL на основе входящего запроса.
//
// Процесс работы функции включает следующие шаги:
// 1. Генерация короткой версии для каждого URL.
// 2. Сохранение всех URL в базу данных.
// 3. Прикрепление сохранённых URL к конкретно
// му пользователю.
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
func (service URLServiceServer) AddURLBatch(ctx context.Context, in *pb.AddURLBatchRequest) (*pb.AddURLBatchResponse, error) {
	var response pb.AddURLBatchResponse
	var requestItems []request.AddURLBatchRequestDTO
	for _, inItem := range in.Items {
		requestItems = append(requestItems, request.AddURLBatchRequestDTO{
			CorrelationID: inItem.CorrelationId,
			OriginalURL:   inItem.OriginalUrl,
		})
	}

	result, err := service.addURLBatchService.AddURLBatch(ctx, requestItems)
	if err != nil {
		switch {
		case errors.Is(err, domainerror.ErrNotValidRequest):
			return nil, status.Errorf(codes.FailedPrecondition, `url deleted`)
		default:
			return nil, status.Errorf(codes.Internal, `internal error`)
		}
	}

	for _, resItem := range result {
		response.Items = append(response.Items, &pb.AddURLBatchResponse_Item{
			CorrelationId: resItem.CorrelationID,
			ShortUrl:      resItem.ShortURL,
		})
	}

	return &response, nil
}

// GetURL получает URL из БД по короткому URL и осуществляет редирект по нему
//
// На вход в URLе приходит сокращенный URL: JRU9a8
// На выход: Оригинальный URL
func (service URLServiceServer) GetURL(ctx context.Context, in *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	var response pb.GetURLResponse

	urlStr, err := service.getURLService.GetURL(in.GetShortUrl())
	if err != nil {
		switch {
		case errors.Is(err, domainerror.ErrURLDeleted):
			return nil, status.Errorf(codes.FailedPrecondition, `url deleted`)
		default:
			return nil, status.Errorf(codes.Internal, `internal error`)
		}
	}

	response.OriginalUrl = urlStr

	return &response, nil
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
func (service URLServiceServer) GetUserURLs(ctx context.Context, in *pb.Empty) (*pb.GetUserURLsResponse, error) {
	var response pb.GetUserURLsResponse

	result, err := service.getUserURLBatchService.GetUserURLs(ctx)
	if err != nil {
		switch {
		case errors.Is(err, domainerror.ErrUserUnauthorized):
			return nil, status.Errorf(codes.Unauthenticated, `user unauthorized`)
		default:
			return nil, status.Errorf(codes.Internal, `internal error`)
		}
	}

	for _, item := range result {
		response.Items = append(response.Items, &pb.GetUserURLsResponse_URL{
			ShortUrl:    item.ShortURL,
			OriginalUrl: item.OriginalURL,
		})
	}
	return &response, nil
}

// DeleteUserURLs удаляет несколько URL, которые создал пользователь.
// Под удалением подразумевается пометка в БД флага is_deleted=true.
// Обработка происходит асинхронно.
//
// На вход принимает массив сокращенных версий URL, созданных пользователем:
// [
//
//	"6qxTVvsy",
//	"RTfd56hn",
//	"Jlfd67ds"
//
// ]
func (service URLServiceServer) DeleteUserURLs(ctx context.Context, in *pb.DeleteUserURLsRequest) (*pb.Empty, error) {
	err := service.deleteUserURLsService.DeleteUserURLs(ctx, in.GetShortUrls())
	if err != nil {
		return nil, status.Errorf(codes.Internal, `internal error`)
	}
	return &pb.Empty{}, nil
}

// DeleteURLBatch удаляет несколько URL.
// На вход принимает сокращенные версии URL:
// [
//
//	"g95W3D",
//	"L7ibuA",
//	"TnZWBr"
//
// ]
func (service URLServiceServer) DeleteURLBatch(ctx context.Context, in *pb.DeleteURLBatchRequest) (*pb.Empty, error) {
	err := service.urlRepository.DeleteURLBatch(ctx, in.GetShortUrls())
	if err != nil {
		return nil, status.Errorf(codes.Internal, `internal error`)
	}
	return &pb.Empty{}, nil
}
