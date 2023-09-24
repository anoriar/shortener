package router

import (
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	addURLHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/shared/response"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Router struct {
	addHandler   *addurlhandler.AddHandler
	getHandler   *geturlhandler.GetHandler
	addHandlerV2 *addURLHandlerV2.AddHandler
	logger       *zap.Logger
}

func NewRouter(addHandler *addurlhandler.AddHandler, getHandler *geturlhandler.GetHandler, addHandlerV2 *addURLHandlerV2.AddHandler, logger *zap.Logger) *Router {
	return &Router{
		addHandler:   addHandler,
		getHandler:   getHandler,
		addHandlerV2: addHandlerV2,
		logger:       logger,
	}
}

func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Use(r.loggerMiddleware)

	router.Post("/", r.addHandler.AddURL)
	router.Get("/{id}", r.getHandler.GetURL)
	router.Post("/api/shorten", r.addHandlerV2.AddURL)

	return router
}

func (r *Router) loggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		timeStart := time.Now()

		lw := response.NewLoggingResponseWriter(w)

		h.ServeHTTP(lw, request)

		duration := time.Since(timeStart)

		responseData := lw.ResponseData()

		//#MENTOR: Целесообразно ли использовать Sugarize? Он удобнее, но напрягает то, что его можно дешугоризовывать динамически.
		r.logger.Info("Request",
			zap.String("uri", request.URL.String()),
			zap.String("method", request.Method),
			zap.String("duration", duration.String()),
			zap.Int("status", responseData.Status()),
			zap.Int("size", responseData.Size()),
		)
	})
}
