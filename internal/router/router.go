package router

import (
	"github.com/anoriar/shortener/internal/handlers/add_url_handler"
	"github.com/anoriar/shortener/internal/handlers/get_url_handler"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	addHandler *add_url_handler.AddHandler
	getHandler *get_url_handler.GetHandler
}

func NewRouter(addHandler *add_url_handler.AddHandler, getHandler *get_url_handler.GetHandler) *Router {
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
