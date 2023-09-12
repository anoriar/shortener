package router

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/handlers/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/geturlhandler"
)

func InitializeRouter(cnf *config.Config) *Router {
	return NewRouter(
		addurlhandler.InitializeAddHandler(cnf),
		geturlhandler.InitializeGetHandler(),
	)
}
