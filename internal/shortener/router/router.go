package router

import (
	"database/sql"
	"fmt"
	"github.com/anoriar/shortener/internal/shortener/config"
	database "github.com/anoriar/shortener/internal/shortener/db"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/ping"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander"
	addURLHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/middleware/compress"
	loggerMiddlewarePkg "github.com/anoriar/shortener/internal/shortener/middleware/logger"
	"github.com/anoriar/shortener/internal/shortener/repository"
	dbURLRepository "github.com/anoriar/shortener/internal/shortener/repository/db"
	"github.com/anoriar/shortener/internal/shortener/repository/file"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/util"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Router struct {
	addHandler         *addurlhandler.AddHandler
	getHandler         *geturlhandler.GetHandler
	addHandlerV2       *addURLHandlerV2.AddHandler
	addURLBatchHandler *addurlbatchhander.AddURLBatchHandler
	pingHandler        *ping.PingHandler
	loggerMiddleware   *loggerMiddlewarePkg.LoggerMiddleware
	compressMiddleware *compress.CompressMiddleware
}

func InitializeRouter(cnf *config.Config, logger *zap.Logger, db *sql.DB) (*Router, error) {
	urlRepository := repository.NewInMemoryURLRepository()

	switch {
	case cnf.DatabaseDSN != "" && db != nil:
		err := database.PrepareDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("database preparing error %s", err)
		}
		urlRepository = dbURLRepository.NewDBURLRepository(db, logger)
	case cnf.FileStoragePath != "":
		urlRepository = file.NewFileURLRepository(cnf.FileStoragePath)
	}

	return NewRouter(
		addurlhandler.NewAddHandler(urlRepository, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), logger, cnf.BaseURL),
		geturlhandler.NewGetHandler(urlRepository, logger),
		addURLHandlerV2.NewAddHandler(urlRepository, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), logger, cnf.BaseURL),
		addurlbatchhander.InitializeAddURLBatchHandler(urlRepository, util.NewKeyGen(), logger, cnf.BaseURL),
		ping.NewPingHandler(db, logger),
		loggerMiddlewarePkg.NewLoggerMiddleware(logger),
		compress.NewCompressMiddleware(),
	), nil
}

func NewRouter(
	addHandler *addurlhandler.AddHandler,
	getHandler *geturlhandler.GetHandler,
	addHandlerV2 *addURLHandlerV2.AddHandler,
	addURLBatchHandler *addurlbatchhander.AddURLBatchHandler,
	pingHandler *ping.PingHandler,
	loggerMiddleware *loggerMiddlewarePkg.LoggerMiddleware,
	compressMiddleware *compress.CompressMiddleware,
) *Router {
	return &Router{
		addHandler:         addHandler,
		getHandler:         getHandler,
		addHandlerV2:       addHandlerV2,
		addURLBatchHandler: addURLBatchHandler,
		pingHandler:        pingHandler,
		loggerMiddleware:   loggerMiddleware,
		compressMiddleware: compressMiddleware,
	}
}

func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Use(r.loggerMiddleware.Log)
	router.Use(r.compressMiddleware.Compress)

	router.Get("/ping", r.pingHandler.Ping)
	router.Post("/", r.addHandler.AddURL)
	router.Get("/{id}", r.getHandler.GetURL)
	router.Post("/api/shorten", r.addHandlerV2.AddURL)
	router.Post("/api/shorten/batch", r.addURLBatchHandler.AddURLBatch)

	return router
}
