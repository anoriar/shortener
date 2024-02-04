package auth

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/util/utilip"
)

// InternalAuthMiddleware missing godoc.
type InternalAuthMiddleware struct {
	conf      *config.Config
	ipService *utilip.IPService
	logger    *zap.Logger
}

// NewInternalAuthMiddleware missing godoc.
func NewInternalAuthMiddleware(conf *config.Config, logger *zap.Logger) *InternalAuthMiddleware {
	return &InternalAuthMiddleware{conf: conf, logger: logger}
}

// InternalAuth missing godoc.
func (iam *InternalAuthMiddleware) InternalAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if iam.conf.TrustedSubnet == "" {

			http.Error(w, "access forbidden", http.StatusForbidden)
			return
		}
		ip, err := iam.ipService.GetIPFromRequest(request)

		if err != nil {
			if errors.Is(err, utilip.ErrIPNotFound) {
				http.Error(w, "access forbidden", http.StatusForbidden)
			} else {
				iam.logger.Error("get short URLs from user error", zap.String("error", err.Error()))
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}

		isIPBelongsToSubnet, err := iam.ipService.IsIPBelongToSubnet(ip, iam.conf.TrustedSubnet)

		if err != nil {
			iam.logger.Error("get short URLs from user error", zap.String("error", err.Error()))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if !isIPBelongsToSubnet {
			http.Error(w, "access forbidden", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, request.WithContext(ctx))
	})

}
