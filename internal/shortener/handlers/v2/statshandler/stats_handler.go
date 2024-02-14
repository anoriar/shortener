package statshandler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/services/stats"
)

// StatsHandler Обработчик получения статистики
type StatsHandler struct {
	statsService stats.StatsServiceInterface
	logger       *zap.Logger
}

// NewStatsHandler missing godoc.
func NewStatsHandler(statsService stats.StatsServiceInterface, logger *zap.Logger) *StatsHandler {
	return &StatsHandler{statsService: statsService, logger: logger}
}

// GetStats missing godoc.
func (handler StatsHandler) GetStats(w http.ResponseWriter, req *http.Request) {
	responseDTO, err := handler.statsService.GetStats(req.Context())
	if err != nil {
		handler.logger.Error("marshal error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	jsonResult, err := json.Marshal(responseDTO)
	if err != nil {
		handler.logger.Error("marshal error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResult)
	if err != nil {
		handler.logger.Error("write response error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
