package middleware

import (
	"github.com/anoriar/shortener/internal/shortener/shared/response"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	logger *zap.Logger
}

func NewLoggerMiddleware(logger *zap.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger: logger}
}

func (lm *LoggerMiddleware) Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		timeStart := time.Now()

		lw := response.NewLoggingResponseWriter(w)

		h.ServeHTTP(lw, request)

		duration := time.Since(timeStart)

		responseData := lw.ResponseData()

		lm.logger.Info("Request",
			zap.String("uri", request.URL.String()),
			zap.String("method", request.Method),
			zap.String("duration", duration.String()),
			zap.Int("status", responseData.Status()),
			zap.Int("size", responseData.Size()),
		)
	})
}
