package router

import (
	"github.com/anoriar/shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	addHandler *handlers.AddHandler
	getHandler *handlers.GetHandler
}

func NewRouter(addHandler *handlers.AddHandler, getHandler *handlers.GetHandler) *Router {
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
