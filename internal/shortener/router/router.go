package router

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	addURLHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/middleware/compress"
	loggerMiddlewarePkg "github.com/anoriar/shortener/internal/shortener/middleware/logger"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/repository/file"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Router struct {
	addHandler         *addurlhandler.AddHandler
	getHandler         *geturlhandler.GetHandler
	addHandlerV2       *addURLHandlerV2.AddHandler
	loggerMiddleware   *loggerMiddlewarePkg.LoggerMiddleware
	compressMiddleware *compress.CompressMiddleware
}

func InitializeRouter(cnf *config.Config, logger *zap.Logger) *Router {
	urlRepository := repository.NewInMemoryURLRepository()
	if cnf.FileStoragePath != "" {
		urlRepository = file.NewFileURLRepository(cnf.FileStoragePath)
	}

	return NewRouter(
		addurlhandler.InitializeAddHandler(cnf, urlRepository),
		geturlhandler.InitializeGetHandler(urlRepository),
		addURLHandlerV2.Initialize(cnf, urlRepository),
		loggerMiddlewarePkg.NewLoggerMiddleware(logger),
		compress.NewCompressMiddleware(),
	)
}

func NewRouter(
	addHandler *addurlhandler.AddHandler,
	getHandler *geturlhandler.GetHandler,
	addHandlerV2 *addURLHandlerV2.AddHandler,
	loggerMiddleware *loggerMiddlewarePkg.LoggerMiddleware,
	compressMiddleware *compress.CompressMiddleware,
) *Router {
	return &Router{
		addHandler:         addHandler,
		getHandler:         getHandler,
		addHandlerV2:       addHandlerV2,
		loggerMiddleware:   loggerMiddleware,
		compressMiddleware: compressMiddleware,
	}
}

func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Use(r.loggerMiddleware.Log)
	router.Use(r.compressMiddleware.Compress)

	router.Post("/", r.addHandler.AddURL)
	router.Get("/{id}", r.getHandler.GetURL)
	router.Post("/api/shorten", r.addHandlerV2.AddURL)

	return router
}
