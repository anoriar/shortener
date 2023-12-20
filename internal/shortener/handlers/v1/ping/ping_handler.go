package ping

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/repository/url"
)

type PingHandler struct {
	urlRepository url.URLRepositoryInterface
	logger        *zap.Logger
}

func NewPingHandler(urlRepository url.URLRepositoryInterface, logger *zap.Logger) *PingHandler {
	return &PingHandler{urlRepository: urlRepository, logger: logger}
}

func (p *PingHandler) Ping(w http.ResponseWriter, req *http.Request) {
	err := p.urlRepository.Ping(req.Context())
	if err != nil {
		p.logger.Error("Storage error", zap.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
