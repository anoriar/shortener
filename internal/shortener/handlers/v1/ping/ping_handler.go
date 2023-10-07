package ping

import (
	"database/sql"
	"github.com/anoriar/shortener/internal/shortener/db"
	"go.uber.org/zap"
	"net/http"
)

type PingHandler struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPingHandler(db *sql.DB, logger *zap.Logger) *PingHandler {
	return &PingHandler{db: db, logger: logger}
}

func (p *PingHandler) Ping(w http.ResponseWriter, req *http.Request) {
	err := db.PingDatabase(p.db)
	if err != nil {
		p.logger.Error("Database error", zap.String("error", err.Error()))
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
