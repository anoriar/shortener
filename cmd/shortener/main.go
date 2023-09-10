package main

import (
	"github.com/anoriar/shortener/internal/config"
	"github.com/anoriar/shortener/internal/di"
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

	ctn, err := di.NewContainer(conf)
	if err != nil {
		panic(err)
	}

	r := ctn.Resolve("router").(*router.Router)

	err = http.ListenAndServe(conf.Host, r.Route())
	if err != nil {
		panic(err)
	}

	err = ctn.Clean()
	if err != nil {
		panic(err)
	}
}
