package router

import (
	"github.com/anoriar/shortener/internal/handlers/addurlhandler"
	"github.com/anoriar/shortener/internal/handlers/geturlhandler"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	addHandler *addurlhandler.AddHandler
	getHandler *geturlhandler.GetHandler
}

func NewRouter(addHandler *addurlhandler.AddHandler, getHandler *geturlhandler.GetHandler) *Router {
	return &Router{
		addHandler: addHandler,
		getHandler: getHandler,
	}
}

func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Post("/", r.addHandler.AddURL)
	router.Get("/{id}", r.getHandler.GetURL)

	return router
}
