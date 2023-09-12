package router

import (
	"github.com/anoriar/shortener/internal/config"
	"github.com/anoriar/shortener/internal/handlers/addurlhandler"
	"github.com/anoriar/shortener/internal/handlers/geturlhandler"
)

func InitializeRouter(cnf *config.Config) *Router {
	return NewRouter(
		addurlhandler.InitializeAddHandler(cnf),
		geturlhandler.InitializeGetHandler(),
	)
}
