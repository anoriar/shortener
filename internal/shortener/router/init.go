package router

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/handlers/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/geturlhandler"
	"github.com/anoriar/shortener/internal/shortener/storage"
)

func InitializeRouter(cnf *config.Config) *Router {
	storage := storage.NewURLStorage()
	return NewRouter(
		addurlhandler.InitializeAddHandler(cnf, storage),
		geturlhandler.InitializeGetHandler(storage),
	)
}
