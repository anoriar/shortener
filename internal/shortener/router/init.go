package router

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	addURLHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/storage"
	"go.uber.org/zap"
)

func InitializeRouter(cnf *config.Config, logger *zap.Logger) *Router {
	storage := storage.NewURLStorage()
	return NewRouter(
		addurlhandler.InitializeAddHandler(cnf, storage),
		geturlhandler.InitializeGetHandler(storage),
		addURLHandlerV2.Initialize(cnf, storage),
		logger,
	)
}
