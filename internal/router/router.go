package router

import (
	"github.com/anoriar/shortener/internal/handlers/addURLHandler"
	"github.com/anoriar/shortener/internal/handlers/getURLHandler"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	addHandler *addURLHandler.AddHandler
	getHandler *getURLHandler.GetHandler
}

func NewRouter(addHandler *addURLHandler.AddHandler, getHandler *getURLHandler.GetHandler) *Router {
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
