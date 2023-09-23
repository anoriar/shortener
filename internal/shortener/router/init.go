package router

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/handlers/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/geturlhandler"
	"github.com/anoriar/shortener/internal/shortener/storage"
	"go.uber.org/zap"
)

func InitializeRouter(cnf *config.Config, logger *zap.Logger) *Router {
	storage := storage.NewURLStorage()
	return NewRouter(
		addurlhandler.InitializeAddHandler(cnf, storage),
		geturlhandler.InitializeGetHandler(storage),
		logger,
	)
}
