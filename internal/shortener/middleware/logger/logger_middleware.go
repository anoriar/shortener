package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/middleware/logger/internal/responsewriter"
)

// LoggerMiddleware missing godoc.
type LoggerMiddleware struct {
	logger *zap.Logger
}

// NewLoggerMiddleware missing godoc.
func NewLoggerMiddleware(logger *zap.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger: logger}
}

// Log missing godoc.
func (lm *LoggerMiddleware) Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		timeStart := time.Now()

		lw := responsewriter.NewLoggingResponseWriter(w)

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
