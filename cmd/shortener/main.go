package main

import (
	"github.com/anoriar/shortener/internal/handlers"
	"github.com/anoriar/shortener/internal/router"
	"github.com/anoriar/shortener/internal/storage"
	"github.com/anoriar/shortener/internal/util"
	"net/http"
)

func main() {
	run()
}

func run() {
	r := router.NewRouter(
		handlers.NewAddHandler(storage.GetInstance(), util.NewKeyGen()),
		handlers.NewGetHandler(storage.GetInstance()),
	)

	http.ListenAndServe("localhost:8080", r.Route())
}
