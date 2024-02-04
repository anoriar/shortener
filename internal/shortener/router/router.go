package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/anoriar/shortener/internal/shortener/handlers/v2/statshandler"

	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/ping"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander"
	addURLHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/deleteurlbatchhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/deleteuserurlshandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/getuserurlshandler"
	"github.com/anoriar/shortener/internal/shortener/middleware/auth"
	"github.com/anoriar/shortener/internal/shortener/middleware/compress"
	loggerMiddlewarePkg "github.com/anoriar/shortener/internal/shortener/middleware/logger"
)

// Router missing godoc.
type Router struct {
	addHandler             *addurlhandler.AddHandler
	getHandler             *geturlhandler.GetHandler
	addHandlerV2           *addURLHandlerV2.AddHandler
	addURLBatchHandler     *addurlbatchhander.AddURLBatchHandler
	getUserURLsHandler     *getuserurlshandler.GetUserURLsHandler
	pingHandler            *ping.PingHandler
	deleteURLBatchHandler  *deleteurlbatchhandler.DeleteURLBatchHandler
	deleteUserURLsHandler  *deleteuserurlshandler.DeleteUserURLsHandler
	statsHandler           *statshandler.StatsHandler
	loggerMiddleware       *loggerMiddlewarePkg.LoggerMiddleware
	compressMiddleware     *compress.CompressMiddleware
	authMiddleware         *auth.AuthMiddleware
	internalAuthMiddleware *auth.InternalAuthMiddleware
}

// NewRouter missing godoc.
func NewRouter(
	addHandler *addurlhandler.AddHandler,
	getHandler *geturlhandler.GetHandler,
	addHandlerV2 *addURLHandlerV2.AddHandler,
	addURLBatchHandler *addurlbatchhander.AddURLBatchHandler,
	getUserURLsHandler *getuserurlshandler.GetUserURLsHandler,
	pingHandler *ping.PingHandler,
	deleteURLBatchHandler *deleteurlbatchhandler.DeleteURLBatchHandler,
	deleteUserURLsHandler *deleteuserurlshandler.DeleteUserURLsHandler,
	statsHandler *statshandler.StatsHandler,
	loggerMiddleware *loggerMiddlewarePkg.LoggerMiddleware,
	compressMiddleware *compress.CompressMiddleware,
	authMiddleware *auth.AuthMiddleware,
	internalAuthMiddleware *auth.InternalAuthMiddleware,
) *Router {
	return &Router{
		addHandler:             addHandler,
		getHandler:             getHandler,
		addHandlerV2:           addHandlerV2,
		addURLBatchHandler:     addURLBatchHandler,
		getUserURLsHandler:     getUserURLsHandler,
		pingHandler:            pingHandler,
		deleteURLBatchHandler:  deleteURLBatchHandler,
		deleteUserURLsHandler:  deleteUserURLsHandler,
		statsHandler:           statsHandler,
		loggerMiddleware:       loggerMiddleware,
		compressMiddleware:     compressMiddleware,
		authMiddleware:         authMiddleware,
		internalAuthMiddleware: internalAuthMiddleware,
	}
}

// Route missing godoc.
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
	router.With(r.authMiddleware.Auth).Get("/api/user/urls", r.getUserURLsHandler.GetUserURLs)
	router.With(r.authMiddleware.Auth).Delete("/api/user/urls", r.deleteUserURLsHandler.DeleteUserURLs)
	router.With(r.internalAuthMiddleware.InternalAuth).Get("/api/internal/stats", r.statsHandler.GetStats)

	return router
}
