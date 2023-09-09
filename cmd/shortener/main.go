package main

import (
	"github.com/anoriar/shortener/internal/config"
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
	conf := config.NewConfig()
	parseFlags(conf)

	r := router.NewRouter(
		handlers.NewAddHandler(storage.GetInstance(), util.NewKeyGen(), conf.BaseURL),
		handlers.NewGetHandler(storage.GetInstance()),
	)

	err := http.ListenAndServe(conf.Host, r.Route())
	if err != nil {
		return
	}
}
