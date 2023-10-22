package router

import (
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/ping"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander"
	addURLHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/deleteurlbatchhandler"
	"github.com/anoriar/shortener/internal/shortener/middleware/auth"
	"github.com/anoriar/shortener/internal/shortener/middleware/compress"
	loggerMiddlewarePkg "github.com/anoriar/shortener/internal/shortener/middleware/logger"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	addHandler            *addurlhandler.AddHandler
	getHandler            *geturlhandler.GetHandler
	addHandlerV2          *addURLHandlerV2.AddHandler
	addURLBatchHandler    *addurlbatchhander.AddURLBatchHandler
	pingHandler           *ping.PingHandler
	deleteURLBatchHandler *deleteurlbatchhandler.DeleteURLBatchHandler
	loggerMiddleware      *loggerMiddlewarePkg.LoggerMiddleware
	compressMiddleware    *compress.CompressMiddleware
	authMiddleware        *auth.AuthMiddleware
}

func NewRouter(
	addHandler *addurlhandler.AddHandler,
	getHandler *geturlhandler.GetHandler,
	addHandlerV2 *addURLHandlerV2.AddHandler,
	addURLBatchHandler *addurlbatchhander.AddURLBatchHandler,
	pingHandler *ping.PingHandler,
	deleteURLBatchHandler *deleteurlbatchhandler.DeleteURLBatchHandler,
	loggerMiddleware *loggerMiddlewarePkg.LoggerMiddleware,
	compressMiddleware *compress.CompressMiddleware,
	authMiddleware *auth.AuthMiddleware,
) *Router {
	return &Router{
		addHandler:            addHandler,
		getHandler:            getHandler,
		addHandlerV2:          addHandlerV2,
		addURLBatchHandler:    addURLBatchHandler,
		pingHandler:           pingHandler,
		deleteURLBatchHandler: deleteURLBatchHandler,
		loggerMiddleware:      loggerMiddleware,
		compressMiddleware:    compressMiddleware,
		authMiddleware:        authMiddleware,
	}
}

func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Use(r.loggerMiddleware.Log)
	router.Use(r.compressMiddleware.Compress)

	router.Get("/ping", r.pingHandler.Ping)
	router.With(r.authMiddleware.Auth).Post("/", r.addHandler.AddURL)
	router.With(r.authMiddleware.Auth).Get("/{id}", r.getHandler.GetURL)
	router.With(r.authMiddleware.Auth).Post("/api/shorten", r.addHandlerV2.AddURL)
	router.With(r.authMiddleware.Auth).Post("/api/shorten/batch", r.addURLBatchHandler.AddURLBatch)
	router.With(r.authMiddleware.Auth).Delete("/api/shorten/batch", r.deleteURLBatchHandler.DeleteURLBatch)

	return router
}
