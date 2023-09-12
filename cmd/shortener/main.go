package main

import (
	"github.com/anoriar/shortener/internal/config"
	"github.com/anoriar/shortener/internal/router"
	"github.com/caarlos0/env/v6"
	"net/http"
)

func main() {
	run()
}

func run() {
	conf := config.NewConfig()
	parseFlags(conf)

	err := env.Parse(conf)
	if err != nil {
		return
	}

	r := router.InitializeRouter(conf)

	err = http.ListenAndServe(conf.Host, r.Route())
	if err != nil {
		panic(err)
	}
}
